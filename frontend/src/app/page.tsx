'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { usePouchDBSync } from '@/hooks/usePouchDBSync'; // 作成したフックをインポート
import { useOccurrenceData } from '@/hooks/useOccurrenceData'; // ★追加: 発生データ取得フック

const DB_NAME = process.env.NEXT_PUBLIC_DB_NAME || 'test_db';

// 型を再定義
interface Workstation {
  workstation_id: number;
  workstation_name: string;
}

// useOccurrenceData から取得するデータの型を再定義
interface Occurrence {
  _id: string;
  created_at: string;
  occurrence_data: {
    individual_id: string;
    sex: string;
  };
  classification_data: {
    class_classification: {
        genus: string;
        species: string;
    }
  }
}

export default function Home() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<any>(null);
  const [workstations, setWorkstations] = useState<Workstation[]>([]); 
  const [currentWS, setCurrentWS] = useState<Workstation | null>(null); 

  // ★ 同期フックを使用
  const syncState = usePouchDBSync(currentWS?.workstation_id || null, user?.user_id || null);
  // ★ 追加: 発生データ取得フックを使用
  const { occurrences, loading: dataLoading, error: dataError } = useOccurrenceData(currentWS?.workstation_id || null);

  useEffect(() => {
    const t = localStorage.getItem('auth_token');
    if (!t) {
      router.push('/login');
      return;
    }
    setToken(t);

    // ユーザー情報の取得
    fetch('/api/users/me', { headers: { Authorization: `Bearer ${t}` } })
      .then(res => res.json())
      .then(data => setUser(data))
      .catch(() => router.push('/login'));

    // ワークステーション一覧取得
    fetch('/api/my-workstations', { headers: { Authorization: `Bearer ${t}` } })
      .then(res => res.json())
      .then((data: Workstation[]) => { 
        setWorkstations(data || []);
        
        if (data && data.length > 0) {
          const savedWSJson = localStorage.getItem('current_workstation');
          let targetWS: Workstation = data[0]; // デフォルトは最新のWS (APIが最新順だと仮定)

          // 修正されたワークステーション選択ロジック (前回の修正を保持)
          if (savedWSJson) {
            try {
              const savedWS = JSON.parse(savedWSJson);
              
              // 1. savedWSが現在の一覧に存在するか確認
              const foundWS = data.find(w => w.workstation_id === savedWS.workstation_id);

              if (foundWS) {
                  // 存在すればそれを優先する 
                  targetWS = foundWS; 
              } 
            } catch (e) {
                console.error("Failed to parse saved workstation:", e);
                // パース失敗時は data[0] (targetWSの初期値) を使う
            }
          }
          
          // 状態とローカルストレージを常に targetWS で更新するのだ
          setCurrentWS(targetWS);
          localStorage.setItem('current_workstation', JSON.stringify(targetWS));
        } else {
            setCurrentWS(null);
            localStorage.removeItem('current_workstation');
        }
      });
  }, [router]);

  const handleWorkstationChange = (ws: Workstation) => {
    setCurrentWS(ws);
    localStorage.setItem('current_workstation', JSON.stringify(ws));
    // ページリロードしてDB再接続等を確実にしても良い
    window.location.reload();
  };

  if (!token) return null;

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow p-4 flex justify-between items-center">
        <h1 className="text-xl font-bold text-black">Web Occurrence</h1>
        <div className="flex items-center gap-4">
          {/* 同期状態の表示 */}
          <span className={`text-xs font-mono px-2 py-1 rounded ${
            syncState === 'active' ? 'bg-green-100 text-green-800' : 
            syncState === 'error' ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-black'
          }`}>
            Sync: {syncState.toUpperCase()}
          </span>

          <select 
            className="border rounded p-1 text-black"
            value={currentWS?.workstation_id || ''}
            onChange={(e) => {
              const ws = workstations.find(w => String(w.workstation_id) === e.target.value);
              if (ws) handleWorkstationChange(ws);
            }}
          >
            {workstations.map(ws => (
              <option key={ws.workstation_id} value={ws.workstation_id}>
                {ws.workstation_name}
              </option>
            ))}
          </select>
          <div className="text-black">{user?.user_name}</div>
        </div>
      </header>

      <main className="p-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* 新規作成カード */}
          <Link href="/create" className="block group">
            <div className="bg-white p-6 rounded-xl shadow-sm border-2 border-dashed border-gray-300 hover:border-green-500 transition-colors h-full flex flex-col items-center justify-center min-h-[200px]">
              <span className="text-4xl text-black group-hover:text-green-500 mb-2">+</span>
              <span className="font-bold text-black group-hover:text-green-600">新規データ登録</span>
            </div>
          </Link>

          {/* ★修正: データ一覧表示エリア */}
          <div className="bg-white p-6 rounded-xl shadow-sm md:col-span-2">
            <h2 className="font-bold mb-4 text-gray-800">最近のデータ</h2>
            {dataError && <p className="text-red-500">エラー: {dataError}</p>}
            {dataLoading ? (
              <p className="text-gray-500">データを読み込み中...</p>
            ) : (
              <OccurrenceList occurrences={occurrences} />
            )}
          </div>
        </div>
      </main>
    </div>
  );
}

// ★追加: 発生データ表示用コンポーネント
function OccurrenceList({ occurrences }: { occurrences: Occurrence[] }) {
    if (occurrences.length === 0) {
        return <p className="text-sm text-gray-500">まだこのワークステーションにはデータが登録されていません。</p>;
    }

    // 最新の5件のみ表示
    const recentOccurrences = occurrences.slice(0, 5);

    return (
        <div className="space-y-3">
            {recentOccurrences.map(occ => (
                <div key={occ._id} className="p-3 border border-gray-100 rounded-lg hover:bg-gray-50 transition">
                    <p className="text-sm font-semibold text-black">
                        {/* 分類情報 (genus, species) を表示 */}
                        {occ.classification_data.class_classification.genus || '不明種'} {occ.classification_data.class_classification.species || '不明'}
                    </p>
                    <p className="text-xs text-gray-500 mt-0.5">
                        ID: {occ._id.substring(0, 8)} | 
                        性別: {occ.occurrence_data.sex} | 
                        日時: {new Date(occ.created_at).toLocaleString()}
                    </p>
                </div>
            ))}
        </div>
    );
}
