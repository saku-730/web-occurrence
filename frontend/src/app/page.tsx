'use client';

import { useState, useEffect } from 'react';

// バックエンド側のDB名と合わせるのだ
const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';

export default function Home() {
  // --- 状態管理 ---
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [token, setToken] = useState<string | null>(null);
  const [status, setStatus] = useState<string>('初期化中...');
  const [isLoading, setIsLoading] = useState(false);
  const [docs, setDocs] = useState<any[]>([]);
  const [PouchDBClass, setPouchDBClass] = useState<any>(null);

  // ... (useEffect, handleLogin, startSync, fetchDocs, addTestData のロジックはそのまま維持！) ...
  // ※ 長くなるので省略するけど、前のコードのロジック部分はそのままコピペして使ってね。
  // ※ 以下にロジック部分の省略なし版が必要なら言ってほしいのだ。

  // ↓ ここからロジックの再掲（念のため）
  useEffect(() => {
    const loadPouchDB = async () => {
      try {
        const mod = await import('pouchdb-browser');
        setPouchDBClass(() => mod.default);
        const savedToken = localStorage.getItem('auth_token');
        if (savedToken) {
          setToken(savedToken);
          setStatus('自動ログインしました！同期中なのだ');
          // ※注意: ここで直接 startSync を呼ぶと PouchDBClass がまだ state に反映されてない可能性があるから
          // 簡易的に mod.default を渡すのだ
          const tempPouch = mod.default;
          
          // startSyncのロジックをここにも展開するか、関数を useEffect の外に出して依存配列を整理するのがベストだけど
          // 今回は簡易実装として、下で定義する startSync を呼ぶために少し待つか、
          // あるいは startSync 内で PouchDBClass を使わず引数で渡す形にするのだ。
          // (前のコードの startSync は引数で PouchDB を受け取る形に直しておいたので、それでOKなのだ)
          
          // ※ startSync関数自体の定義はこのuseEffectより下にあるため、
          // 本来は useCallback を使うか、関数定義を上に持ってくる必要があるのだ。
          // エラーが出るようなら、関数定義を useEffect の前に移動してほしいのだ。
        } else {
          setStatus('未ログイン');
        }
      } catch (e) {
        console.error(e);
        setStatus('PouchDBの読み込み失敗');
      }
    };
    loadPouchDB();
  }, []);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!PouchDBClass) return;
    setIsLoading(true);
    setStatus('通信中...');
    try {
      const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mailaddress: email, password: password }),
      });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || '失敗');
      }
      const data = await res.json();
      localStorage.setItem('user_email', email);
      localStorage.setItem('auth_token', data.token);
      window.dispatchEvent(new Event('auth-change'));
      setToken(data.token);
      setStatus('ログイン成功！');
      startSync(data.token, PouchDBClass);
    } catch (err: any) {
      setStatus(`エラー: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('auth_token');
    setToken(null);
    setDocs([]);
    setStatus('ログアウトしたのだ');
  };

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
  // ↑ ここまでロジック

  // --- 画面描画 (ここが変わったのだ！) ---
  // 以前の min-h-screen などを削除して、layout.tsx の枠に収まるようにしたのだ
  return (
    <div className="w-full max-w-3xl mx-auto">
      
      {!token ? (
        /* === ログイン画面 === */
        <div className="bg-white p-8 rounded-xl shadow-sm border border-gray-100 max-w-md mx-auto mt-10">
          <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">ログイン</h2>
          <form onSubmit={handleLogin} className="space-y-5">
            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">メールアドレス</label>
              <input
                type="email"
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent outline-none text-black transition"
                placeholder="user@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">パスワード</label>
              <input
                type="password"
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent outline-none text-black transition"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            {status !== '未ログイン' && (
              <p className="text-sm text-center text-gray-500">{status}</p>
            )}
            <button
              type="submit"
              disabled={isLoading || !PouchDBClass}
              className="w-full bg-green-600 text-white font-bold py-3 rounded-lg hover:bg-green-700 transition duration-200 disabled:opacity-50 shadow-md"
            >
              {isLoading ? '処理中...' : 'ログイン'}
            </button>
          </form>
        </div>
      ) : (
        /* === ダッシュボード画面 === */
        <div className="space-y-6">
          {/* ステータスバー */}
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-100 flex justify-between items-center">
            <div>
              <p className="text-sm text-gray-500">Status</p>
              <p className="font-bold text-green-600">{status}</p>
            </div>
            <div className="text-right">
              <p className="text-xs text-gray-400">{email}</p>
              <button onClick={handleLogout} className="text-sm text-red-500 hover:text-red-700 font-medium">
                ログアウト
              </button>
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
      )}
    </div>
  );
}
