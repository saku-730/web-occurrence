import { useEffect, useState } from 'react';

// PouchDBの動的インポート用ヘルパー
const getPouchDB = async () => {
  const mod = await import('pouchdb-browser');
  // 型エラー回避のためのキャスト
  return ((mod.default && typeof mod.default === 'function') ? mod.default : mod) as unknown as any;
};

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';
// CouchDBのURL (環境変数またはデフォルト)
// 注意: ブラウザからアクセスする場合、CouchDB側でCORS設定が必要なのだ
const REMOTE_COUCHDB_URL = process.env.NEXT_PUBLIC_COUCHDB_URL || 'http://localhost:5984';

export function usePouchDBSync(workstationId: number | string | null) {
  const [syncState, setSyncState] = useState<'stopped' | 'active' | 'paused' | 'error'>('stopped');

  useEffect(() => {
    if (!workstationId) return;

    let localDB: any;
    let remoteDB: any;
    let syncHandler: any;

    const startSync = async () => {
      const PouchDB = await getPouchDB();
      const dbName = `${DB_NAME}_ws_${workstationId}`;
      
      // ローカルDB
      localDB = new PouchDB(dbName);
      
      // リモートDB (CouchDB)
      // 認証が必要な場合は { auth: { username: '...', password: '...' } } などを第2引数に追加する
      const remoteUrl = `${REMOTE_COUCHDB_URL}/${dbName}`;
      remoteDB = new PouchDB(remoteUrl);

      console.log(`[Sync] Starting sync for ${dbName} to ${remoteUrl}`);

      // 双方向同期 (Live)
      syncHandler = localDB.sync(remoteDB, {
        live: true,
        retry: true
      }).on('change', (info: any) => {
        console.log('[Sync] Change detected:', info);
        setSyncState('active');
      }).on('paused', (err: any) => {
        console.log('[Sync] Paused (Offline?)', err);
        setSyncState('paused');
      }).on('active', () => {
        console.log('[Sync] Active');
        setSyncState('active');
      }).on('error', (err: any) => {
        console.error('[Sync] Error:', err);
        setSyncState('error');
      });
    };

    startSync();

    // クリーンアップ
    return () => {
      if (syncHandler) syncHandler.cancel();
      if (localDB) localDB.close();
      if (remoteDB) remoteDB.close();
    };
  }, [workstationId]);

  return syncState;
}
