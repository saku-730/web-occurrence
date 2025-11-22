'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function CreateWorkstationPage() {
  const router = useRouter();
  const [wsName, setWsName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const token = localStorage.getItem('auth_token');
      if (!token) throw new Error('ログイン情報の取得に失敗しました');

      const res = await fetch('/api/workstation/create', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ workstation_name: wsName }),
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || '作成に失敗しました');
      }

      // 作成成功！トップページへ戻って同期を開始するのだ
      router.push('/');

    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-4">
      <div className="bg-white p-8 rounded-xl shadow-sm border border-gray-100 w-full max-w-md">
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">ワークステーション作成</h2>
        
        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="block text-sm font-semibold text-gray-600 mb-2">名称</label>
            <input
              type="text"
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 outline-none text-black transition"
              placeholder="例: つくばミミズ調査2025"
              value={wsName}
              onChange={(e) => setWsName(e.target.value)}
            />
            <p className="text-xs text-gray-400 mt-1">わかりやすい名前をつけてください</p>
          </div>

          {error && (
            <p className="text-sm text-center text-red-500 bg-red-50 p-2 rounded">
              {error}
            </p>
          )}

          <div className="pt-2">
            <button
              type="submit"
              disabled={isLoading || !wsName}
              className="w-full bg-green-600 text-white font-bold py-3 rounded-lg hover:bg-green-700 transition duration-200 disabled:opacity-50 shadow-md"
            >
              {isLoading ? '作成中...' : '作成して開始'}
            </button>
          </div>
          
          <div className="text-center">
            <button 
              type="button"
              onClick={() => router.back()}
              className="text-sm text-gray-500 hover:text-gray-700"
            >
              戻る
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
