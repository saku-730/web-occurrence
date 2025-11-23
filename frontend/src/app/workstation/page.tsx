'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

interface Workstation {
  workstation_id: number;
  workstation_name: string;
}

export default function WorkstationListPage() {
  const router = useRouter();
  const [workstations, setWorkstations] = useState<Workstation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchWorkstations = async () => {
      try {
        const token = localStorage.getItem('auth_token');
        if (!token) {
          router.push('/login');
          return;
        }

        // ★修正: /api/workstations -> /api/my-workstations に変更
        // バックエンドの定義と一致させ、自分の所属するワークステーションを取得するようにしたのだ
        const res = await fetch('/api/my-workstations', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (!res.ok) throw new Error('データの取得に失敗しました');
        
        const data = await res.json();
        setWorkstations(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchWorkstations();
  }, [router]);

  const handleSelect = (ws: Workstation) => {
    // 選択したワークステーションをローカルストレージに保存
    localStorage.setItem('current_workstation', JSON.stringify(ws));
    
    // トップページへ移動（そこでPouchDBの初期化が行われる想定）
    router.push('/');
  };

  if (loading) return <div className="p-8 text-center">読み込み中...</div>;

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <h1 className="text-2xl font-bold text-gray-800">ワークステーションを選択</h1>
        
        {error && <p className="text-red-500">{error}</p>}

        {/* ワークステーション一覧 */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {workstations.map((ws) => (
            <button
              key={ws.workstation_id}
              onClick={() => handleSelect(ws)}
              className="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:border-green-500 hover:shadow-md transition text-left"
            >
              <div className="h-10 w-10 bg-green-100 text-green-600 rounded-full flex items-center justify-center mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-6 h-6">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M20.25 14.15v4.25c0 1.094-.787 2.036-1.872 2.18-2.087.277-4.216.42-6.378.42s-4.291-.143-6.378-.42c-1.085-.144-1.872-1.086-1.872-2.18v-4.25m16.5 0a2.18 2.18 0 00.75-1.661V8.706c0-1.081-.768-2.015-1.837-2.175a48.114 48.114 0 00-3.413-.387m4.5 8.053c-.211 1.806-.777 3.587-1.643 5.26m0 0c-.9 1.734-2.786 3.19-5.752 4.674m5.752-4.674V12.75" />
                </svg>
              </div>
              <h2 className="font-bold text-lg text-gray-800">{ws.workstation_name}</h2>
              <p className="text-xs text-gray-500 mt-1">ID: {ws.workstation_id}</p>
            </button>
          ))}
        </div>

        {/* アクションボタンエリア */}
        <div className="border-t border-gray-200 pt-8">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">新しいワークステーション</h2>
          <div className="flex gap-4">
            <Link 
              href="/workstation/create" 
              className="flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition shadow-sm font-medium"
            >
              <span>新規作成</span>
            </Link>
            
            <button 
              onClick={() => alert('申請機能は準備中なのだ')}
              className="flex items-center gap-2 px-6 py-3 bg-white text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition shadow-sm font-medium"
            >
              <span>申請する</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
