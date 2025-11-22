const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';
const MASTER_DOC_ID = '_local/master_data';

// APIのレスポンス型定義
interface MasterApiResponse {
  languages: any[];
  file_types: any[];
  file_extensions: any[];
  user_roles: any[];
  workstation_users: any[];
}

export async function fetchAndSaveMasterData(jwt: string) {
  console.log('マスターデータの同期を開始します...');
  
  // SSRエラー回避のため動的インポートを使用
  // トップレベルで import PouchDB from 'pouchdb-browser' するとNext.jsのビルドで落ちる可能性があるのだ
  const pouchModule = await import('pouchdb-browser');
  const PouchDB = pouchModule.default || pouchModule;
  
  const db = new PouchDB(DB_NAME);

  try {
    // 1. バックエンドからマスターデータを取得
    // ★修正: バックエンドの router.go に合わせて GET メソッドを使うのだ
    const res = await fetch('/api/master-data', {
      method: 'GET', 
      headers: {
        'Authorization': `Bearer ${jwt}`,
        'Content-Type': 'application/json',
      },
    });

    if (!res.ok) {
      throw new Error(`API Error: ${res.status} ${res.statusText}`);
    }

    const apiData: MasterApiResponse = await res.json();

    // 2. 現在のローカルドキュメントを取得（_rev取得のため）
    let currentDoc: any = {};
    try {
      currentDoc = await db.get(MASTER_DOC_ID);
    } catch (err: any) {
      if (err.status !== 404) throw err;
      // 存在しない場合は新規作成
    }

    // 3. データに変更があるか簡易チェック
    if (currentDoc.data && JSON.stringify(currentDoc.data) === JSON.stringify(apiData)) {
      console.log('マスターデータに変更はありませんでした。');
      return;
    }

    // 4. 保存するデータの構築
    const newDoc = {
      _id: MASTER_DOC_ID,
      _rev: currentDoc._rev, // 更新時は必須
      data: apiData,
      updated_at: new Date().toISOString(),
    };

    // 5. PouchDBに保存
    await db.put(newDoc);
    console.log('マスターデータを更新しました！');

  } catch (err) {
    console.error('マスターデータの同期に失敗しました:', err);
    // マスター同期の失敗はメイン処理を止めないため、ログ出力のみにとどめる
  }
}
