'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';

export default function CreateOccurrencePage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [currentWS, setCurrentWS] = useState<any>(null);
  const [token, setToken] = useState<string | null>(null);
  
  // フォーム状態：分類階級をフルサポート
  const [formData, setFormData] = useState({
    kingdom: 'Animalia',
    phylum: '',
    class: '',
    order: '',
    family: '',
    genus: '',
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
      const uuid = () => crypto.randomUUID();
      const userId = 'user_saku'; // TODO: 実際のログインユーザーIDを使う

      const occurrenceId = uuid();
      const classificationId = uuid();
      const placeId = uuid();
      
      const doc = {
        _id: `occ_${occurrenceId}`,
        type: 'occurrence',
        workstation_id: String(currentWS.workstation_id),
        created_by_user_id: userId,
        project_id: null,
        created_at: new Date(formData.date).toISOString(),
        timezone: '+09:00',
        language_id: '1',

        occurrence_data: {
          individual_id: formData.individual_id,
          lifestage: 'adult',
          sex: formData.sex,
          body_length: null,
          note: formData.note,
        },

        classification_data: {
          classification_id: classificationId,
          class_classification: {
            kingdom: formData.kingdom,
            phylum: formData.phylum,
            class: formData.class,
            order: formData.order,
            family: formData.family,
            genus: formData.genus,
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

        identifications: [],
        specimens: [],
        observations: [],
        attachments: []
      };

      // ▼▼▼ 修正箇所: 型定義エラーを確実に回避するのだ ▼▼▼
      const mod = await import('pouchdb-browser');
      // as unknown as any で強制的に型を無視させるのだ
      const PouchDB = ((mod.default && typeof mod.default === 'function') ? mod.default : mod) as unknown as any;
      
      const localDBName = `${DB_NAME}_ws_${currentWS.workstation_id}`;
      const db = new PouchDB(localDBName);
      // ▲▲▲ 修正ここまで ▲▲▲
      
      await db.put(doc);

      console.log('Saved to local PouchDB:', doc);
      alert('保存しました！');
      
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
      <div className="max-w-2xl mx-auto bg-white p-8 rounded-xl shadow-sm text-black">
        <h1 className="text-2xl font-bold mb-6">新規データ登録</h1>
        <p className="text-sm text-gray-500 mb-4">
          Workstation: {currentWS.workstation_name} (ID: {currentWS.workstation_id})
        </p>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* 分類情報の入力エリア */}
          <div className="space-y-2 border p-4 rounded bg-gray-50">
            <h2 className="font-bold text-lg mb-2">分類情報 (Classification)</h2>
            
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-bold mb-1">Kingdom (界)</label>
                <input name="kingdom" value={formData.kingdom} onChange={handleChange} className="w-full border p-2 rounded" placeholder="Animalia" />
              </div>
              <div>
                <label className="block text-sm font-bold mb-1">Phylum (門)</label>
                <input name="phylum" value={formData.phylum} onChange={handleChange} className="w-full border p-2 rounded" placeholder="Arthropoda" />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-bold mb-1">Class (綱)</label>
                <input name="class" value={formData.class} onChange={handleChange} className="w-full border p-2 rounded" placeholder="Insecta" />
              </div>
              <div>
                <label className="block text-sm font-bold mb-1">Order (目)</label>
                <input name="order" value={formData.order} onChange={handleChange} className="w-full border p-2 rounded" placeholder="Lepidoptera" />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-bold mb-1">Family (科)</label>
                <input name="family" value={formData.family} onChange={handleChange} className="w-full border p-2 rounded" placeholder="Papilionidae" />
              </div>
              <div>
                <label className="block text-sm font-bold mb-1">Genus (属)</label>
                <input name="genus" value={formData.genus} onChange={handleChange} className="w-full border p-2 rounded" placeholder="Papilio" />
              </div>
            </div>

            <div>
              <label className="block text-sm font-bold mb-1">Species (種名)</label>
              <input name="species" value={formData.species} onChange={handleChange} required className="w-full border p-2 rounded" placeholder="xuthus (種小名のみ、または学名全体)" />
            </div>
          </div>

          <hr className="my-4" />

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
