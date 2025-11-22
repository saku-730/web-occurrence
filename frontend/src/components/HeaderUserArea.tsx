'use client';

import { useState, useEffect } from 'react';

export default function HeaderUserArea() {
  // --- 状態管理 ---
  const [isOnline, setIsOnline] = useState(true);
  const [userEmail, setUserEmail] = useState<string | null>(null);
  const [mounted, setMounted] = useState(false);

  // --- 初期化とイベント監視 ---
  useEffect(() => {
    setMounted(true);

    // 1. オンライン状態の監視
    setIsOnline(navigator.onLine);
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    // 2. ユーザー情報の取得
    const email = localStorage.getItem('user_email');
    setUserEmail(email);

    // 3. ストレージの変更を監視 (ログイン/ログアウト時の表示更新用)
    const handleStorageChange = () => {
      setUserEmail(localStorage.getItem('user_email'));
    };
    // カスタムイベント 'auth-change' を監視
    window.addEventListener('auth-change', handleStorageChange);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      window.removeEventListener('auth-change', handleStorageChange);
    };
  }, []);

  // --- ログアウト処理 ---
  const handleLogout = () => {
    if (!confirm('ログアウトしますか？')) return;
    
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_email');
    
    // イベントを発火して他コンポーネントに通知
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
      {userEmail ? (
        <div className="flex items-center gap-4 pl-4 border-l border-gray-200">
          <span className="text-gray-600 hidden sm:inline-block">
            {userEmail}
          </span>
          <button 
            onClick={handleLogout}
            className="text-gray-500 hover:text-red-600 transition font-medium"
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
