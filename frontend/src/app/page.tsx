'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { usePouchDBSync } from '@/hooks/usePouchDBSync'; // 作成したフックをインポート

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';

export default function Home() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<any>(null);
  const [workstations, setWorkstations] = useState<any[]>([]);
  const [currentWS, setCurrentWS] = useState<any>(null);

  // ★ 同期フックを使用
  const syncState = usePouchDBSync(currentWS?.workstation_id || null);

  useEffect(() => {
    const t = localStorage.getItem('auth_token');
    if (!t) {
      router.push('/login');
      return;
    }
    setToken(t);

    // ユーザー情報のデコード (簡易実装: 本来はjwt-decode等を使う)
    // ログイン時にlocalStorageにuser_idなどを保存しておくのがベスト
    // ここではAPIから取得する流れを想定
    fetch('/api/users/me', { headers: { Authorization: `Bearer ${t}` } })
      .then(res => res.json())
      .then(data => setUser(data))
      .catch(() => router.push('/login'));

    // ワークステーション一覧取得
    fetch('/api/my-workstations', { headers: { Authorization: `Bearer ${t}` } })
      .then(res => res.json())
      .then(data => {
        setWorkstations(data || []);
        // 簡易的に最初のWSを選択状態にする
        if (data && data.length > 0) {
          const savedWS = localStorage.getItem('current_workstation');
          if (savedWS) {
            setCurrentWS(JSON.parse(savedWS));
          } else {
            setCurrentWS(data[0]);
            localStorage.setItem('current_workstation', JSON.stringify(data[0]));
          }
        }
      });
  }, [router]);

  const handleWorkstationChange = (ws: any) => {
    setCurrentWS(ws);
    localStorage.setItem('current_workstation', JSON.stringify(ws));
    // ページリロードしてDB再接続等を確実にしても良い
    window.location.reload();
  };

  if (!token) return null;

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow p-4 flex justify-between items-center">
        <h1 className="text-xl font-bold">Web Occurrence</h1>
        <div className="flex items-center gap-4">
          {/* 同期状態の表示 */}
          <span className={`text-xs font-mono px-2 py-1 rounded ${
            syncState === 'active' ? 'bg-green-100 text-green-800' : 
            syncState === 'error' ? 'bg-red-100 text-red-800' : 'bg-gray-100'
          }`}>
            Sync: {syncState.toUpperCase()}
          </span>

          <select 
            className="border rounded p-1"
            value={currentWS?.workstation_id || ''}
            onChange={(e) => {
              const ws = workstations.find(w => String(w.workstation_id) === e.target.value);
              if (ws) handleWorkstationChange(ws);
            }}
          >
            {workstations.map(ws => (
              <option key={ws.workstation_id} value={ws.workstation_id}>
                {ws.workstation_name}
              </option>
            ))}
          </select>
          <div>{user?.user_name}</div>
        </div>
      </header>

      <main className="p-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* 新規作成カード */}
          <Link href="/create" className="block group">
            <div className="bg-white p-6 rounded-xl shadow-sm border-2 border-dashed border-gray-300 hover:border-green-500 transition-colors h-full flex flex-col items-center justify-center min-h-[200px]">
              <span className="text-4xl text-gray-400 group-hover:text-green-500 mb-2">+</span>
              <span className="font-bold text-gray-500 group-hover:text-green-600">新規データ登録</span>
            </div>
          </Link>

          {/* データ一覧表示エリア (PouchDBから取得して表示する実装が必要) */}
          {/* 今回は省略するが、usePouchDBSyncと同じ要領で allDocs を取得して表示する */}
          <div className="bg-white p-6 rounded-xl shadow-sm">
            <h2 className="font-bold mb-2">最近のデータ</h2>
            <p className="text-sm text-gray-500">ここに登録済みデータが表示されます。</p>
          </div>
        </div>
      </main>
    </div>
  );
}
