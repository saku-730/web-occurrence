'use client';

import { useState, useEffect } from 'react';
import PouchDB from 'pouchdb-browser';

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';
const MASTER_DOC_ID = '_local/master_data';

export default function MasterDataPage() {
  const [masterData, setMasterData] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      try {
        const db = new PouchDB(DB_NAME);
        // _local/master_data ドキュメントを取得
        const doc: any = await db.get(MASTER_DOC_ID);
        setMasterData(doc.data);
      } catch (err: any) {
        console.error(err);
        if (err.status === 404) {
          setError('マスターデータが見つかりません。一度トップページに戻って同期を待ってください。');
        } else {
          setError('データの読み込みに失敗しました。');
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
        <span className="text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded">Local Cache</span>
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
