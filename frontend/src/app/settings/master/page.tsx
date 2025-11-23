'use client';

import { useState, useEffect } from 'react';
// PouchDBは動的インポートで使用するため、ここではimportしない

const MASTER_DOC_ID = '_local/master_data';

export default function MasterDataPage() {
  const [masterData, setMasterData] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      try {
        // 1. ローカルストレージから現在選択中のワークステーション情報を取得
        // Next.jsなどのSSR環境を考慮し、windowオブジェクトの存在確認を行うのが安全だが
        // useEffect内であればクライアントサイドでの実行が保証されるのでそのままアクセス可能
        const wsJson = localStorage.getItem('current_workstation');
        
        if (!wsJson) {
          setError('ワークステーションが選択されていません。トップページまたはワークステーション選択画面に戻ってください。');
          setLoading(false);
          return;
        }

        const currentWorkstation = JSON.parse(wsJson);
        const wsId = currentWorkstation.workstation_id;

        if (!wsId) {
          setError('ワークステーションIDが無効です。再度ワークステーションを選択し直してください。');
          setLoading(false);
          return;
        }

        // 2. データベース名を動的に構築 (例: db_ws_1)
        // バックエンドや同期フックの命名規則と一致させる必要がある
        const dbName = `db_ws_${wsId}`;
        console.log(`[MasterPage] Loading from DB: ${dbName}`);

        // 3. PouchDBを動的にインポート
        const pouchModule = await import('pouchdb-browser');
        const PouchDB = pouchModule.default || pouchModule;
        
        // 4. 正しいDB名でインスタンス化
        const db = new PouchDB(dbName);

        // 5. _local/master_data ドキュメントを取得
        const doc: any = await db.get(MASTER_DOC_ID);
        
        if (!doc.data) {
          throw new Error('マスターデータの形式が不正です (dataプロパティがありません)');
        }

        setMasterData(doc.data);

      } catch (err: any) {
        console.error('[MasterPage] Error:', err);
        if (err.status === 404) {
          setError(`マスターデータが見つかりません。同期がまだ完了していない可能性があります。(DB: db_ws_${JSON.parse(localStorage.getItem('current_workstation') || '{}').workstation_id})`);
        } else {
          setError(`データの読み込みに失敗しました: ${err.message || err}`);
        }
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  if (loading) return <div className="p-8 text-center text-gray-500">読み込み中...</div>;
  if (error) return <div className="p-8 text-center text-red-500">{error}</div>;
  if (!masterData) return <div className="p-8 text-center">データがありません</div>;

  return (
    <div className="w-full max-w-6xl mx-auto space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">マスターデータ一覧</h1>
        <div className="flex items-center gap-2">
           <span className="text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded">
             WS ID: {JSON.parse(localStorage.getItem('current_workstation') || '{}').workstation_id}
           </span>
           <span className="text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded">Local Cache</span>
        </div>
      </div>

      {/* --- 言語一覧 --- */}
      <MasterTable 
        title="言語 (Languages)" 
        data={masterData.languages} 
        columns={[
          { key: 'language_id', label: 'ID' },
          { key: 'language_short', label: '略称' },
          { key: 'language_common', label: '名称' },
        ]} 
      />

      {/* --- ファイル種別 --- */}
      <MasterTable 
        title="ファイル種別 (File Types)" 
        data={masterData.file_types} 
        columns={[
          { key: 'file_type_id', label: 'ID' },
          { key: 'type_name', label: '種別名' },
        ]} 
      />

      {/* --- 拡張子 --- */}
      <MasterTable 
        title="拡張子 (File Extensions)" 
        data={masterData.file_extensions} 
        columns={[
          { key: 'extension_id', label: 'ID' },
          { key: 'extension_text', label: '拡張子' },
          { key: 'file_type_id', label: '種別ID' },
        ]} 
      />

      {/* --- ユーザーロール --- */}
      <MasterTable 
        title="ユーザーロール (User Roles)" 
        data={masterData.user_roles} 
        columns={[
          { key: 'role_id', label: 'ID' },
          { key: 'role_name', label: 'ロール名' },
        ]} 
      />

      {/* --- ワークステーションユーザー --- */}
      <MasterTable 
        title="ワークステーションユーザー (WS Users)" 
        data={masterData.workstation_users} 
        columns={[
          { key: 'user_id', label: 'User ID' },
          { key: 'display_name', label: '表示名' },
        ]} 
      />
    </div>
  );
}

// テーブル表示用のサブコンポーネント
function MasterTable({ title, data, columns }: { title: string, data: any[], columns: { key: string, label: string }[] }) {
  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="px-6 py-3 bg-gray-50 border-b border-gray-100">
        <h2 className="font-bold text-gray-700 text-sm">{title}</h2>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-sm text-left text-gray-600">
          <thead className="text-xs text-gray-700 uppercase bg-white border-b border-gray-100">
            <tr>
              {columns.map((col) => (
                <th key={col.key} className="px-6 py-3 font-medium">
                  {col.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {(!data || data.length === 0) ? (
              <tr>
                <td colSpan={columns.length} className="px-6 py-4 text-center text-gray-400">
                  データなし
                </td>
              </tr>
            ) : (
              data.map((row: any, idx: number) => (
                <tr key={idx} className="bg-white border-b border-gray-50 hover:bg-gray-50">
                  {columns.map((col) => (
                    <td key={col.key} className="px-6 py-3">
                      {row[col.key]}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
