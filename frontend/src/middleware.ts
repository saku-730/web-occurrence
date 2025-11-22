import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  // 1. Cookieから認証トークンを取得するのだ
  const token = request.cookies.get('auth_token');

  // 2. 現在のパスを取得
  const path = request.nextUrl.pathname;

  // 3. 認証不要なパス（ログイン、登録、静的ファイル系）を定義
  const isPublicPath = path === '/login' || path === '/register';

  // --- リダイレクトロジック ---

  // A. トークンがなくて、かつ保護されたページ（/login, /register以外）にアクセスした場合
  // -> ログイン画面へ強制転送
  if (!token && !isPublicPath) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // B. すでにトークンがあるのに、ログイン画面や登録画面にアクセスした場合
  // -> トップページ（ダッシュボード）へ転送（親切設計）
  if (token && isPublicPath) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  // 何もなければそのまま通過させる
  return NextResponse.next();
}

// ミドルウェアを適用するパスの設定
export const config = {
  // api, _next, favicon などを除外して、それ以外すべてに適用
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
