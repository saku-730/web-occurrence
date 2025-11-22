'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function Register() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [status, setStatus] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setStatus('登録処理中...');

    try {
      // バックエンドの登録APIを叩く
      const res = await fetch('/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          mailaddress: email, 
          password: password 
        }),
      });

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || '登録に失敗しました');
      }

      // 成功したらアラートを出してログイン画面へ
      alert('登録が完了したのだ！ログインしてほしいのだ。');
      router.push('/');

    } catch (err: any) {
      setStatus(`エラー: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-4">
      <div className="bg-white p-8 rounded-xl shadow-sm border border-gray-100 w-full max-w-md">
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">新規ユーザー登録</h2>
        
        <form onSubmit={handleRegister} className="space-y-5">
          <div>
            <label className="block text-sm font-semibold text-gray-600 mb-2">メールアドレス</label>
            <input
              type="email"
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 outline-none text-black transition"
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
              placeholder="8文字以上"
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
            className="w-full bg-blue-600 text-white font-bold py-3 rounded-lg hover:bg-blue-700 transition duration-200 disabled:opacity-50 shadow-md"
          >
            {isLoading ? '登録中...' : '登録する'}
          </button>
        </form>

        <div className="mt-6 text-center">
          <Link href="/" className="text-sm text-green-600 hover:underline">
            すでにアカウントをお持ちの方はこちら（ログイン）
          </Link>
        </div>
      </div>
    </div>
  );
}
