'use client';

import Link from 'next/link';

export default function NewWorkstationPage() {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-2xl text-center space-y-8">
        <h1 className="text-3xl font-bold text-gray-800">ようこそ！</h1>
        <p className="text-gray-600">
          まずは作業場所（ワークステーション）を設定しましょう。<br />
          新しく作成するか、チームのワークステーションに参加するか選んでください。
        </p>

        <div className="grid md:grid-cols-2 gap-6">
          {/* 作成カード */}
          <Link 
            href="/workstation/create"
            className="bg-white p-8 rounded-xl shadow-sm border border-gray-200 hover:border-green-500 hover:shadow-md transition group"
          >
            <div className="h-12 w-12 bg-green-100 text-green-600 rounded-full flex items-center justify-center mx-auto mb-4 group-hover:bg-green-600 group-hover:text-white transition">
              {/* プラスアイコン */}
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-6 h-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
              </svg>
            </div>
            <h2 className="text-xl font-bold text-gray-800 mb-2">新しく作成する</h2>
            <p className="text-sm text-gray-500">
              自分用のデータベースを作成して、<br />ゼロからデータを記録します。
            </p>
          </Link>

          {/* 参加カード */}
          <button 
            onClick={() => alert('まだ実装されてないのだ！')}
            className="bg-white p-8 rounded-xl shadow-sm border border-gray-200 hover:border-blue-500 hover:shadow-md transition group text-left w-full"
          >
            <div className="h-12 w-12 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 group-hover:bg-blue-600 group-hover:text-white transition">
              {/* 参加アイコン */}
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-6 h-6">
                <path strokeLinecap="round" strokeLinejoin="round" d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
              </svg>
            </div>
            <h2 className="text-xl font-bold text-center text-gray-800 mb-2">既存に参加する</h2>
            <p className="text-sm text-center text-gray-500">
              他のメンバーが作成した<br />ワークステーションに参加します。
            </p>
          </button>
        </div>
      </div>
    </div>
  );
}
