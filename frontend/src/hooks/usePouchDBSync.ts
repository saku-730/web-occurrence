import { useEffect, useState } from 'react';

// PouchDBの動的インポート用ヘルパー
const getPouchDB = async () => {
  const mod = await import('pouchdb-browser');
  // 型エラー回避のためのキャスト
  return ((mod.default && typeof mod.default === 'function') ? mod.default : mod) as unknown as any;
};

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';
const BACKEND_BASE_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8080';

export function usePouchDBSync(workstationId: number | string | null) {
  const [syncState, setSyncState] = useState<'stopped' | 'active' | 'paused' | 'error'>('stopped');

  useEffect(() => {
    if (!workstationId || typeof window === 'undefined') return;

    const token = localStorage.getItem('auth_token');
    if (!token) {
      console.warn('[Sync] トークンがないため同期を停止します');
      return;
    }

    let localDB: any;
    let remoteDB: any;
    let syncHandler: any;

    const startSync = async () => {
      try {
        const PouchDB = await getPouchDB();
        const dbName = `${DB_NAME}_ws_${workstationId}`;
        
        // 1. ローカルDB
        localDB = new PouchDB(dbName);
        
        // 2. リモートDB (Go Proxy)
        const remoteUrl = `${BACKEND_BASE_URL}/api/couchdb/${dbName}`;
        console.log(`[Sync] Connecting to proxy: ${remoteUrl}`);

        // ★重要: fetchをオーバーライドして、強制的にヘッダーを注入するのだ！
        const remoteOpts = {
          skip_setup: true,
          fetch: function (url: string, opts: any) {
            // opts.headers が Headers オブジェクトか単純なオブジェクトか確認してセット
            if (!opts.headers) {
                opts.headers = new Headers();
            } else if (!(opts.headers instanceof Headers)) {
                opts.headers = new Headers(opts.headers);
            }
            
            // ここでトークンをセット！
            opts.headers.set('Authorization', `Bearer ${token}`);
            
            // デバッグ用ログ（本番では消してもいいのだ）
            // console.log('[Sync Fetch]', url, opts.headers.get('Authorization'));

            return PouchDB.fetch(url, opts);
          }
        };

        remoteDB = new PouchDB(remoteUrl, remoteOpts);

        // 3. 同期開始
        syncHandler = localDB.sync(remoteDB, {
          live: true,
          retry: true
        })
        .on('change', (info: any) => {
          console.log('[Sync] Change:', info);
          setSyncState('active');
        })
        .on('paused', (err: any) => {
          console.log('[Sync] Paused:', err);
          setSyncState('paused');
        })
        .on('active', () => {
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

    return () => {
      if (syncHandler) syncHandler.cancel();
      if (localDB) localDB.close();
      if (remoteDB) remoteDB.close();
    };
  }, [workstationId]);

  return syncState;
}
