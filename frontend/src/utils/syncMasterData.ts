const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';
const MASTER_DOC_ID = '_local/master_data';

// APIのレスポンス型定義
interface MasterApiResponse {
  workstation_id: number;
  languages: any[];
  file_types: any[];
  file_extensions: any[];
  user_roles: any[];
  workstation_users: any[];
}

// 引数に PouchDBClass を追加するのだ
export async function fetchAndSaveMasterData(jwt: string, PouchDBClass: any) {
  console.log('マスターデータの同期を開始します...');
  
  // ★修正: 受け取ったものが関数(コンストラクタ)でなければ .default を使う
  const DBClass = (typeof PouchDBClass === 'function') 
    ? PouchDBClass 
    : (PouchDBClass.default || PouchDBClass);

  // インスタンス化
  const db = new DBClass(DB_NAME);

  try {
    // 1. バックエンドからマスターデータを取得
    const res = await fetch('/api/master-data', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${jwt}`,
        'Content-Type': 'application/json',
      },
    });

    if (!res.ok) {
      throw new Error(`API Error: ${res.statusText}`);
    }

    const apiData: MasterApiResponse = await res.json();

    // 2. 現在のローカルドキュメントを取得（_rev取得のため）
    let currentDoc: any = {};
    try {
      currentDoc = await db.get(MASTER_DOC_ID);
    } catch (err: any) {
      if (err.status !== 404) throw err;
      // 存在しない場合は新規作成なので何もしない
    }

    // 3. 保存するデータの構築
    const newDoc = {
      _id: MASTER_DOC_ID,
      _rev: currentDoc._rev,
      data: apiData,
      updated_at: new Date().toISOString(),
    };

    // 4. データに変更があるか簡易チェック
    if (currentDoc.data && JSON.stringify(currentDoc.data) === JSON.stringify(apiData)) {
      console.log('マスターデータに変更はありませんでした。');
      return;
    }

    // 5. PouchDBに保存
    await db.put(newDoc);
    console.log('マスターデータを更新しました！ WorkstationID:', apiData.workstation_id);

  } catch (err) {
    console.error('マスターデータの同期に失敗しました:', err);
  }
}
