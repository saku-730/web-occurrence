import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import Link from 'next/link'
import HeaderUserArea from '@/components/HeaderUserArea' // ★ここを追加

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Web Occurrence',
  description: 'Biological Occurrence Data Management',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body className={`${inter.className} bg-gray-50 text-gray-900 min-h-screen flex flex-col`}>
        
        {/* === ヘッダーエリア === */}
        <header className="bg-white shadow-sm sticky top-0 z-50">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            
            {/* 上段: ロゴエリア + 右端のユーザーエリア */}
            <div className="flex items-center justify-between py-4">
              {/* 左: ロゴ */}
              <Link href="/" className="flex items-center gap-3 hover:opacity-80 transition">
                <div className="h-10 w-10 bg-green-500 rounded-md flex items-center justify-center text-white font-bold text-xs shadow-sm">
                  WO
                </div>
                <span className="text-xl font-bold text-gray-800 tracking-tight">
                  Web Occurrence
                </span>
              </Link>

              {/* 右: ユーザーエリア（ネット状況・メアド・ログアウト） */}
              {/* ★ここに配置したのだ！ */}
              <HeaderUserArea />
            </div>

            {/* 下段: トップバー (ナビゲーション) */}
            <nav className="border-t border-gray-100 py-3 overflow-x-auto scrollbar-hide">
              <ul className="flex space-x-8 text-sm font-medium">
                <li>
                  <Link href="/" className="text-gray-600 hover:text-green-600 transition flex items-center gap-1">
                    ホーム
                  </Link>
                </li>
                <li>
                  <Link href="/search" className="text-gray-600 hover:text-green-600 transition">
                    データ検索
                  </Link>
                </li>
                <li>
                  <Link href="/create" className="text-gray-600 hover:text-green-600 transition">
                    新規登録
                  </Link>
                </li>
                <li>
                  <Link href="/map" className="text-gray-600 hover:text-green-600 transition">
                    地図表示
                  </Link>
                </li>
                <li>
                  <Link href="/settings" className="text-gray-600 hover:text-green-600 transition">
                    設定
                  </Link>
                </li>
              </ul>
            </nav>
          </div>
        </header>

        {/* === メインコンテンツ === */}
        <main className="flex-grow w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {children}
        </main>

        {/* === フッター === */}
        <footer className="bg-white border-t border-gray-200 py-6 mt-auto">
          <div className="max-w-7xl mx-auto px-4 text-center text-gray-500 text-xs">
            &copy; 2025 Web Occurrence System. All rights reserved.
          </div>
        </footer>

      </body>
    </html>
  )
}
