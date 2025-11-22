'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
// マスターデータ同期はトップページで行うので、ここではimport削除またはコメントアウト
// import { fetchAndSaveMasterData } from '@/utils/syncMasterData';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [status, setStatus] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
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
        throw new Error(err.error || 'ログイン失敗');
      }

      const data = await res.json();

      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('user_email', email);
      document.cookie = `auth_token=${data.token}; path=/; max-age=86400; SameSite=Lax`;

      window.dispatchEvent(new Event('auth-change'));
      setStatus('ログイン成功！');
      
      // ★修正: ワークステーション選択画面へ移動するのだ
      router.push('/workstation/new');

    } catch (err: any) {
      setStatus(`エラー: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-4">
      <div className="bg-white p-8 rounded-xl shadow-sm border border-gray-100 w-full max-w-md">
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">ログイン</h2>
        
        <form onSubmit={handleLogin} className="space-y-5">
          <div>
            <label className="block text-sm font-semibold text-gray-600 mb-2">メールアドレス</label>
            <input
              type="email"
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 outline-none text-black transition"
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
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 outline-none text-black transition"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          {status && (
            <p className={`text-sm text-center ${status.includes('エラー') ? 'text-red-500' : 'text-blue-500'}`}>
              {status}
            </p>
          )}
          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-green-600 text-white font-bold py-3 rounded-lg hover:bg-green-700 transition duration-200 disabled:opacity-50 shadow-md"
          >
            {isLoading ? '処理中...' : 'ログイン'}
          </button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-sm text-gray-500">アカウントをお持ちでないですか？</p>
          <Link href="/register" className="text-sm text-blue-600 hover:underline font-medium">
            新規ユーザー登録はこちら
          </Link>
        </div>
      </div>
    </div>
  );
}
