import { useEffect, useState } from 'react';

// PouchDBの動的インポート用ヘルパー
const getPouchDB = async () => {
  const mod = await import('pouchdb-browser');
  // 型エラー回避のためのキャスト
  return ((mod.default && typeof mod.default === 'function') ? mod.default : mod) as unknown as any;
};

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';

// バックエンドのURL（Goサーバー）
// Next.jsのプロキシを通さず直接Goサーバー(8080)を叩く想定ですが、
// Next.jsの設定で /api をGoに飛ばしている場合はそのURLに合わせてください。
// ここではCORS設定済みのGoサーバーを直接指定します。
const BACKEND_BASE_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8080';

export function usePouchDBSync(workstationId: number | string | null) {
  const [syncState, setSyncState] = useState<'stopped' | 'active' | 'paused' | 'error'>('stopped');

  useEffect(() => {
    // ワークステーションIDがない、またはブラウザ環境でない場合は何もしない
    if (!workstationId || typeof window === 'undefined') return;

    // 認証トークンの取得
    const token = localStorage.getItem('auth_token');
    if (!token) {
      console.warn('[Sync] No auth token found. Sync suspended.');
      return;
    }

    let localDB: any;
    let remoteDB: any;
    let syncHandler: any;

    const startSync = async () => {
      try {
        const PouchDB = await getPouchDB();
        const dbName = `${DB_NAME}_ws_${workstationId}`;
        
        // 1. ローカルDB (IndexedDB) の初期化
        localDB = new PouchDB(dbName);
        
        // 2. リモートDB (Go Proxy) の設定
        // URL構造: http://localhost:8080/api/couchdb/{dbName}
        // Go側のルーター設定: apiProtected.Any("/couchdb/*path", ...) に対応
        const remoteUrl = `${BACKEND_BASE_URL}/api/couchdb/${dbName}`;

        console.log(`[Sync] Starting sync for ${dbName} via proxy: ${remoteUrl}`);

        // PouchDBのオプション設定
        // skip_setup: true -> DB作成はput時に任せる（プロキシ経由でのDB作成権限依存）
        // ajax.headers -> Goのミドルウェア認証用のBearerトークンを設定
        const remoteOpts = {
          skip_setup: true,
          ajax: {
            headers: {
              'Authorization': `Bearer ${token}`,
              // 必要に応じて他のヘッダーも追加可能
            },
            timeout: 60000 // タイムアウト設定（必要に応じて調整）
          }
        };

        remoteDB = new PouchDB(remoteUrl, remoteOpts);

        // 3. 双方向同期 (Live Replication) の開始
        syncHandler = localDB.sync(remoteDB, {
          live: true,   // 変更をリアルタイムで監視
          retry: true,  // 接続切断時に再試行
          batch_size: 50 // 一度に同期するドキュメント数（パフォーマンス調整）
        })
        .on('change', (info: any) => {
          console.log('[Sync] Change detected:', info);
          setSyncState('active');
        })
        .on('paused', (err: any) => {
          // 接続切れや待機中（ネットワークエラー等もここで検知されることがある）
          console.log('[Sync] Paused (Waiting for changes or offline):', err);
          setSyncState('paused');
        })
        .on('active', () => {
          // 同期再開
          console.log('[Sync] Active');
          setSyncState('active');
        })
        .on('error', (err: any) => {
          console.error('[Sync] Error:', err);
          setSyncState('error');
        });

      } catch (err) {
        console.error('[Sync] Init Error:', err);
        setSyncState('error');
      }
    };

    startSync();

    // クリーンアップ処理
    return () => {
      if (syncHandler) syncHandler.cancel();
      // close()は非同期だがクリーンアップではawaitできないためそのまま呼ぶ
      if (localDB) localDB.close();
      if (remoteDB) remoteDB.close();
    };
  }, [workstationId]);

  return syncState;
}
