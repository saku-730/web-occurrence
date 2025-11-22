'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
// 作成した同期関数をインポート
import { fetchAndSaveMasterData } from '@/utils/syncMasterData';

// ★ DB名を動的に変えるため定数ではなく関数や変数で扱う必要があるけど、
// 今回は単一DB構成で進めているので、DB名は固定で、フィルタリングやバリデーションで制御する想定にするのだ。
const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';

export default function Home() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [currentWS, setCurrentWS] = useState<any>(null);
  const [status, setStatus] = useState<string>('初期化中...');
  const [docs, setDocs] = useState<any[]>([]);
  const [PouchDBClass, setPouchDBClass] = useState<any>(null);

  useEffect(() => {
    const loadPouchDB = async () => {
      try {
        // PouchDBを動的インポート
        const mod = await import('pouchdb-browser');
        // ★修正: mod.default が関数ならそれを、そうでなければ mod を使う
        const PouchDB = (mod.default && typeof mod.default === 'function') ? mod.default : mod;
        
        setPouchDBClass(() => PouchDB);
        
        const savedToken = localStorage.getItem('auth_token');
        const savedWS = localStorage.getItem('current_workstation');

        if (!savedToken) {
          router.push('/login');
          return;
        }

        // ★ ワークステーションが選択されていなければ選択画面へ
        if (!savedWS) {
          router.push('/workstation');
          return;
        }

        setToken(savedToken);
        setCurrentWS(JSON.parse(savedWS));
        
        // ▼ 追加: マスターデータを取得・保存するのだ！
        setStatus('マスターデータ同期中...');
        // ここでAPIを叩いてIndexedDBに保存する処理を実行
        // 第2引数に動的インポートした PouchDB クラスを渡すのだ！
        await fetchAndSaveMasterData(savedToken, PouchDB);
        
        setStatus('同期準備完了');
        
        // メインデータの同期開始
        startSync(savedToken, PouchDB);

      } catch (e) {
        console.error(e);
        setStatus('PouchDBの読み込み失敗');
      }
    };
    loadPouchDB();
  }, [router]);

  // --- PouchDB 同期処理 ---
  const startSync = (jwt: string, PouchDB: any) => {
    // ★ 安全策: PouchDBがクラスか確認
    const DBClass = (typeof PouchDB === 'function') ? PouchDB : (PouchDB.default || PouchDB);

    // ★ ここで DB名を変えることで「ダウンロードしてそのDBをいじる」を実現できる。
    // ローカルDB名を `test_db_ws_{id}` のように分けるのが一番安全なのだ。
    const ws = JSON.parse(localStorage.getItem('current_workstation') || '{}');
    if (!ws.workstation_id) {
        console.error("Workstation IDが見つかりません");
        return;
    }

    const localDBName = `${DB_NAME}_ws_${ws.workstation_id}`; // ワークステーションごとにローカルDBを分ける
    const localDB = new DBClass(localDBName);
    
    // リモートは1つの巨大なDB (occurrence) なのでそのまま
    // (※ 本来はフィルタリングが必要)
    const remoteDB = new DBClass(`/api/couchdb/${DB_NAME}`, {
      fetch: (url: string, opts: any) => {
        opts.headers.set('Authorization', `Bearer ${jwt}`);
        return DBClass.fetch(url, opts);
      },
    });

    // ★ 本来はここで filter オプションを使って、サーバーから
    // 「このワークステーションのデータだけ」をプルするように設定するのだ。
    // 今回はフィルタ実装までは含まれていないので、全データ同期になる点に注意なのだ。
    localDB.sync(remoteDB, { 
      live: true, 
      retry: true 
    })
      .on('change', () => fetchDocs(localDB))
      .on('error', (err: any) => console.error(err));
      
    fetchDocs(localDB);
  };

  const fetchDocs = async (db: any) => {
    const res = await db.allDocs({ include_docs: true, descending: true });
    setDocs(res.rows.map((row: any) => row.doc));
  };

  if (!token || !currentWS) {
    return <div className="min-h-screen flex items-center justify-center text-gray-500">読み込み中...</div>;
  }

  return (
    <div className="w-full max-w-3xl mx-auto">
      <div className="space-y-6">
        {/* ヘッダー：現在のワークステーション表示 */}
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-bold text-gray-800">
            {currentWS.workstation_name} <span className="text-sm font-normal text-gray-500">(ID: {currentWS.workstation_id})</span>
          </h2>
          <button 
            onClick={() => {
              localStorage.removeItem('current_workstation');
              router.push('/workstation');
            }}
            className="text-sm text-blue-600 hover:underline"
          >
            切替
          </button>
        </div>

        <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-100 flex justify-between items-center">
          <div>
            <p className="text-sm text-gray-500">Status</p>
            <p className="font-bold text-green-600">{status}</p>
          </div>
        </div>

        {/* データ追加ボタンエリアを削除したのだ */}

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
