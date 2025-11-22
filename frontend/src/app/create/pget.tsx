'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';

export default function CreateOccurrencePage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [currentWS, setCurrentWS] = useState<any>(null);
  const [token, setToken] = useState<string | null>(null);
  
  // フォーム状態（簡易版）
  const [formData, setFormData] = useState({
    kingdom: 'Animalia',
    species: '',
    individual_id: '',
    sex: 'unknown',
    date: new Date().toISOString().slice(0, 16), // YYYY-MM-DDTHH:mm
    latitude: '',
    longitude: '',
    note: ''
  });

  useEffect(() => {
    const t = localStorage.getItem('auth_token');
    const ws = localStorage.getItem('current_workstation');
    if (!t || !ws) {
      router.push('/');
      return;
    }
    setToken(t);
    setCurrentWS(JSON.parse(ws));
  }, [router]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      // UUID生成ヘルパー
      const uuid = () => crypto.randomUUID();

      // 現在のユーザーID (簡易的にLocalStorageから取得できればいいが、今回は固定値かTokenからデコードが必要)
      // 本番ではJWTからデコードするべきだが、ここでは仮に "user_dummy" とする（バックエンドで補正するか、本来はログイン時にUser情報を保持すべき）
      const userId = 'user_saku'; // TODO: 実際のログインユーザーIDを使う

      // 1. オカレンスデータの構築
      const occurrenceId = uuid();
      const classificationId = uuid();
      const placeId = uuid();
      
      const doc = {
        _id: `occ_${occurrenceId}`, // PouchDBのID
        type: 'occurrence',
        workstation_id: String(currentWS.workstation_id), // Stringに変換
        created_by_user_id: userId,
        project_id: null,
        created_at: new Date(formData.date).toISOString(),
        timezone: '+09:00',
        language_id: '1',

        occurrence_data: {
          individual_id: formData.individual_id,
          lifestage: 'adult', // 簡易化
          sex: formData.sex,
          body_length: null,
          note: formData.note,
        },

        classification_data: {
          classification_id: classificationId,
          class_classification: {
            kingdom: formData.kingdom,
            species: formData.species
          }
        },

        place_data: {
          place_id: placeId,
          place_name_id: null,
          coordinates: (formData.latitude && formData.longitude) 
            ? { type: 'Point', coordinates: [Number(formData.longitude), Number(formData.latitude)] } 
            : null,
          accuracy: null,
          class_place_name: null
        },

        // 配列データ（今回は空で作成）
        identifications: [],
        specimens: [],
        observations: [],
        attachments: []
      };

      // 2. PouchDBへの保存
      const mod = await import('pouchdb-browser');
      const PouchDB = (mod.default && typeof mod.default === 'function') ? mod.default : mod;
      
      const localDBName = `${DB_NAME}_ws_${currentWS.workstation_id}`;
      const db = new PouchDB(localDBName);
      
      await db.put(doc);

      console.log('Saved to local PouchDB:', doc);
      alert('保存しました！');
      
      // トップページへ戻って同期させる
      router.push('/');

    } catch (err: any) {
      console.error(err);
      alert('保存に失敗しました: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  if (!currentWS) return <div>Loading...</div>;

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-2xl mx-auto bg-white p-8 rounded-xl shadow-sm">
        <h1 className="text-2xl font-bold mb-6">新規データ登録</h1>
        <p className="text-sm text-gray-500 mb-4">
          Workstation: {currentWS.workstation_name} (ID: {currentWS.workstation_id})
        </p>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-bold mb-1">Kingdom (界)</label>
              <input name="kingdom" value={formData.kingdom} onChange={handleChange} className="w-full border p-2 rounded" />
            </div>
            <div>
              <label className="block text-sm font-bold mb-1">Species (種名)</label>
              <input name="species" value={formData.species} onChange={handleChange} required className="w-full border p-2 rounded" placeholder="例: Papilio xuthus" />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-bold mb-1">個体ID</label>
              <input name="individual_id" value={formData.individual_id} onChange={handleChange} className="w-full border p-2 rounded" />
            </div>
            <div>
              <label className="block text-sm font-bold mb-1">性別</label>
              <select name="sex" value={formData.sex} onChange={handleChange} className="w-full border p-2 rounded">
                <option value="unknown">不明</option>
                <option value="male">オス</option>
                <option value="female">メス</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-bold mb-1">日時</label>
            <input type="datetime-local" name="date" value={formData.date} onChange={handleChange} className="w-full border p-2 rounded" />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-bold mb-1">緯度</label>
              <input type="number" step="any" name="latitude" value={formData.latitude} onChange={handleChange} className="w-full border p-2 rounded" placeholder="例: 36.2" />
            </div>
            <div>
              <label className="block text-sm font-bold mb-1">経度</label>
              <input type="number" step="any" name="longitude" value={formData.longitude} onChange={handleChange} className="w-full border p-2 rounded" placeholder="例: 140.1" />
            </div>
          </div>

          <div>
            <label className="block text-sm font-bold mb-1">備考</label>
            <textarea name="note" value={formData.note} onChange={handleChange} className="w-full border p-2 rounded h-24" />
          </div>

          <button type="submit" disabled={loading} className="w-full bg-green-600 text-white font-bold py-3 rounded hover:bg-green-700 disabled:opacity-50">
            {loading ? '保存中...' : '保存する'}
          </button>
        </form>
      </div>
    </div>
  );
}
