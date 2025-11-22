'use client';

import { useState, useEffect } from 'react';
// ログアウト時にCookieも消すためにimportが必要なのだ
import { useRouter } from 'next/navigation';

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';

export default function Home() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [status, setStatus] = useState<string>('初期化中...');
  const [docs, setDocs] = useState<any[]>([]);
  const [PouchDBClass, setPouchDBClass] = useState<any>(null);

  // --- 初期化 ---
  useEffect(() => {
    const loadPouchDB = async () => {
      try {
        const mod = await import('pouchdb-browser');
        setPouchDBClass(() => mod.default);
        
        const savedToken = localStorage.getItem('auth_token');
        if (savedToken) {
          setToken(savedToken);
          setStatus('同期準備完了');
          startSync(savedToken, mod.default);
        } else {
          // Middlewareで弾かれるはずだけど、念のため
          router.push('/login');
        }
      } catch (e) {
        console.error(e);
        setStatus('PouchDBの読み込み失敗');
      }
    };
    loadPouchDB();
  }, [router]);

  // --- PouchDB 同期処理 ---
  const startSync = (jwt: string, PouchDB: any) => {
    const localDB = new PouchDB(DB_NAME);
    const remoteDB = new PouchDB(`/api/couchdb/${DB_NAME}`, {
      fetch: (url: string, opts: any) => {
        opts.headers.set('Authorization', `Bearer ${jwt}`);
        return PouchDB.fetch(url, opts);
      },
    });
    localDB.sync(remoteDB, { live: true, retry: true })
      .on('change', () => fetchDocs(localDB))
      .on('error', (err: any) => console.error(err));
    fetchDocs(localDB);
  };

  const fetchDocs = async (db: any) => {
    const res = await db.allDocs({ include_docs: true, descending: true });
    setDocs(res.rows.map((row: any) => row.doc));
  };

  const addTestData = async () => {
    if (!PouchDBClass) return;
    const db = new PouchDBClass(DB_NAME);
    await db.post({
      type: 'occurrence',
      title: 'New Data',
      created_at: new Date().toISOString(),
      workstation_id: 'ws-01',
      created_by_user_id: '16'
    });
  };

  // --- 画面描画 ---
  // ログインフォームの分岐を削除して、ダッシュボードだけにするのだ
  if (!token) {
    return <div className="min-h-screen flex items-center justify-center text-gray-500">読み込み中...</div>;
  }

  return (
    <div className="w-full max-w-3xl mx-auto">
      <div className="space-y-6">
        {/* ステータスバー */}
        <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-100 flex justify-between items-center">
          <div>
            <p className="text-sm text-gray-500">Status</p>
            <p className="font-bold text-green-600">{status}</p>
          </div>
        </div>

        {/* アクションエリア */}
        <div className="flex justify-end">
          <button
            onClick={addTestData}
            className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition shadow-sm font-medium"
          >
            + データを追加
          </button>
        </div>

        {/* データ一覧 */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-100 bg-gray-50">
            <h3 className="font-bold text-gray-700">データ一覧 ({docs.length})</h3>
          </div>
          <ul className="divide-y divide-gray-100 max-h-[500px] overflow-y-auto">
            {docs.length === 0 ? (
              <li className="p-8 text-center text-gray-400">データがありません</li>
            ) : (
              docs.map((doc: any) => (
                <li key={doc._id} className="p-4 hover:bg-gray-50 transition">
                  <pre className="text-xs text-gray-600 whitespace-pre-wrap break-all font-mono">
                    {JSON.stringify(doc, null, 2)}
                  </pre>
                </li>
              ))
            )}
          </ul>
        </div>
      </div>
    </div>
  );
}
