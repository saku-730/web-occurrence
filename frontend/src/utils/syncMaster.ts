// バックエンドからマスターデータを取得してIndexedDBに保存する共通関数なのだ
export const fetchAndSaveMasterData = async (
  jwt: string, 
  PouchDB: any, 
  dbName: string
) => {
  try {
    console.log("マスターデータを取得中...");
    
    // 1. GoサーバーのAPIを叩く
    const res = await fetch('/api/master-data', {
      headers: { 'Authorization': `Bearer ${jwt}` }
    });

    if (!res.ok) {
      console.warn('マスターデータの取得に失敗したのだ:', res.statusText);
      return;
    }

    const masterData = await res.json();

    // 2. ローカルDB (IndexedDB) に保存 (Upsert)
    const localDB = new PouchDB(dbName);
    const docId = 'master_data'; // 固定ID

    try {
      // 既存のデータがあれば _rev を取得して上書き
      const existingDoc = await localDB.get(docId);
      await localDB.put({
        ...existingDoc, // _id と _rev を維持
        type: 'master_data', // バリデーション用
        data: masterData     // 最新データで上書き
      });
      console.log("マスターデータを更新したのだ！");
    } catch (err: any) {
      if (err.name === 'not_found') {
        // 新規作成
        await localDB.put({
          _id: docId,
          type: 'master_data',
          data: masterData
        });
        console.log("マスターデータを新規保存したのだ！");
      } else {
        throw err;
      }
    }
  } catch (e) {
    console.error("マスターデータ同期エラー:", e);
    // マスターデータの失敗でアプリを止める必要はないのでエラーは握りつぶすのだ
  }
};
