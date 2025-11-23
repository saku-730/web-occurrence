// frontend/src/hooks/useOccurrenceData.ts

import { useState, useEffect } from 'react';

// PouchDBの動的インポート用ヘルパー (usePouchDBSync.ts と同様)
const getPouchDB = async () => {
  const mod = await import('pouchdb-browser');
  return ((mod.default && typeof mod.default === 'function') ? mod.default : mod) as unknown as any;
};

// 発生データ（occurrence）の簡易型定義
interface OccurrenceDoc {
  _id: string;
  _rev: string;
  type: 'occurrence';
  workstation_id: string;
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
  // その他のデータは省略
}

const isOccurrenceDoc = (doc: any): doc is OccurrenceDoc => {
    // _local/ や _design/ ではない、かつ type が 'occurrence' のドキュメントのみを返すのだ
    return doc.type === 'occurrence' && !doc._id.startsWith('_design/') && !doc._id.startsWith('_local/');
}


export function useOccurrenceData(workstationId: number | string | null) {
  const [occurrences, setOccurrences] = useState<OccurrenceDoc[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!workstationId || typeof window === 'undefined') {
      setOccurrences([]);
      setLoading(false);
      return;
    }

    let db: any;
    let changesHandler: any;

    const dbName = `db_ws_${workstationId}`; // 同期フックと同じDB名を使用するのだ

    const fetchData = async () => {
      try {
        const PouchDB = await getPouchDB();
        db = new PouchDB(dbName);

        // 1. データを全て取得するのだ
        const result = await db.allDocs({
          include_docs: true,
          descending: true, // 最新のものを先頭にするために降順にする
        });

        // 2. 取得したドキュメントから発生データのみをフィルタリングするのだ
        const occurrenceList = result.rows
          .map((row: any) => row.doc)
          .filter(isOccurrenceDoc);

        setOccurrences(occurrenceList);
        setLoading(false);

        // 3. ライブリスナーを設定して変更を監視するのだ
        changesHandler = db.changes({
          live: true,
          since: 'now',
          include_docs: true,
        }).on('change', (change: any) => {
          // 変更があったらデータを再取得する（簡易実装）
          console.log('[OccurrenceData] Change detected, refetching.');
          fetchData(); // 再帰的に呼ぶが、changesHandlerの重複登録はしない
        }).on('error', (err: any) => {
          console.error('[OccurrenceData] Changes error:', err);
          setError('データのリアルタイム更新に失敗しました。');
        });

      } catch (err: any) {
        console.error('[OccurrenceData] Fetch Error:', err);
        setError(`データの取得に失敗しました: ${err.message}`);
        setLoading(false);
      }
    };

    fetchData();

    return () => {
      // クリーンアップ: ライブリスナーとDB接続を解除するのだ
      if (changesHandler) changesHandler.cancel();
    };
  }, [workstationId]);

  return { occurrences, loading, error };
}
