'use client';

import { useState, useEffect } from 'react';

export default function HeaderUserArea() {
  // --- 状態管理 ---
  const [isOnline, setIsOnline] = useState(true);
  // emailだけでなく、ログイン状態そのものを管理するフラグを持つ
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [userEmail, setUserEmail] = useState<string | null>(null);
  const [mounted, setMounted] = useState(false);

  // --- 初期化とイベント監視 ---
  useEffect(() => {
    setMounted(true);
    setIsOnline(navigator.onLine);

    // 1. オンライン状態の監視
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);
    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    // 2. ユーザー情報の更新関数
    const updateUserInfo = () => {
      const token = localStorage.getItem('auth_token');
      const email = localStorage.getItem('user_email');
      
      // トークンがあればログイン済みとみなす！
      setIsLoggedIn(!!token);
      setUserEmail(email);
      
      console.log('Header update:', { token: !!token, email }); // デバッグ用ログ
    };

    // 初回実行
    updateUserInfo();

    // 3. イベントリスナー登録
    window.addEventListener('auth-change', updateUserInfo);
    window.addEventListener('storage', updateUserInfo); // 別タブでの変更も検知

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      window.removeEventListener('auth-change', updateUserInfo);
      window.removeEventListener('storage', updateUserInfo);
    };
  }, []);

  // --- ログアウト処理 ---
  const handleLogout = () => {
    if (!confirm('ログアウトしますか？')) return;
    
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_email');
    
    // イベントを発火して通知
    window.dispatchEvent(new Event('auth-change'));
    
    // ページをリロードしてログイン画面に戻す
    window.location.reload();
  };

  // マウント前は何も表示しない
  if (!mounted) return null;

  return (
    <div className="flex items-center gap-4 text-sm">
      {/* ネット接続状況 */}
      <div className={`flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium border ${
        isOnline 
          ? 'bg-green-50 text-green-700 border-green-200' 
          : 'bg-red-50 text-red-700 border-red-200'
      }`}>
        <span className={`h-2 w-2 rounded-full ${isOnline ? 'bg-green-500' : 'bg-red-500'}`} />
        {isOnline ? 'Online' : 'Offline'}
      </div>

      {/* ログインユーザー情報 & ログアウト */}
      {isLoggedIn ? (
        <div className="flex items-center gap-4 pl-4 border-l border-gray-200">
          {/* メールアドレスがあれば表示、なければ 'User' と表示 */}
          <span className="text-gray-600 hidden sm:inline-block font-medium">
            {userEmail || 'User'}
          </span>
          <button 
            onClick={handleLogout}
            className="text-gray-500 hover:text-red-600 transition font-medium text-xs sm:text-sm"
          >
            ログアウト
          </button>
        </div>
      ) : (
        <div className="text-gray-400 text-xs">未ログイン</div>
      )}
    </div>
  );
}
