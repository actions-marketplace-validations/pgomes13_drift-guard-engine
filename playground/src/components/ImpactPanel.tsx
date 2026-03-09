"use client";

import { useState } from "react";
import { DiffResult, Hit } from "@/lib/types";

interface Props {
  diff: DiffResult;
  apiUrl: string;
}

const PLACEHOLDER = `// Paste your service code here to find references to breaking changes.
// Example (Go):
func DeleteUser(id string) error {
    return client.Delete("/users/" + id)
}`;

export default function ImpactPanel({ diff, apiUrl }: Props) {
  const [code, setCode]         = useState("");
  const [filename, setFilename] = useState("service.go");
  const [loading, setLoading]   = useState(false);
  const [hits, setHits]         = useState<Hit[] | null>(null);
  const [error, setError]       = useState<string | null>(null);

  const breakingCount = diff.summary.breaking;

  async function handleScan() {
    setLoading(true);
    setHits(null);
    setError(null);
    try {
      const res = await fetch(`${apiUrl}/api/impact`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ diff, code, filename }),
      });
      const data = await res.json();
      if (!res.ok) {
        setError(data.error ?? "Unknown error");
        return;
      }
      setHits(data as Hit[]);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Network error");
    } finally {
      setLoading(false);
    }
  }

  // Group hits by change_path for display
  const groups = hits
    ? hits.reduce<Record<string, Hit[]>>((acc, h) => {
        (acc[h.change_path] ??= []).push(h);
        return acc;
      }, {})
    : {};

  const distinctFiles = hits ? new Set(hits.map((h) => h.file)).size : 0;

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center gap-3 flex-wrap">
        <span className="text-sm font-semibold text-gray-500">Impact Analysis</span>
        <span className="px-3 py-1 rounded-full text-xs font-semibold bg-red-50 text-red-700 border border-red-300">
          {breakingCount} breaking change{breakingCount !== 1 ? "s" : ""} to scan for
        </span>
      </div>

      {/* Filename + code input */}
      <div className="border border-gray-200 rounded-lg overflow-hidden shadow-sm">
        <div className="flex items-center gap-3 bg-gray-50 border-b border-gray-200 px-4 py-2">
          <span className="text-xs font-medium text-gray-500 shrink-0">Filename</span>
          <input
            type="text"
            value={filename}
            onChange={(e) => setFilename(e.target.value)}
            placeholder="service.go"
            className="flex-1 text-xs font-mono bg-white border border-gray-200 rounded px-2 py-1 text-gray-700 focus:outline-none focus:border-indigo-400 max-w-[200px]"
          />
          <span className="text-xs text-gray-400 ml-auto">
            Paste a file from your service to find references to breaking changes
          </span>
        </div>
        <textarea
          value={code}
          onChange={(e) => setCode(e.target.value)}
          placeholder={PLACEHOLDER}
          rows={12}
          className="w-full font-mono text-sm text-gray-800 bg-white px-4 py-3 resize-y focus:outline-none"
          spellCheck={false}
        />
      </div>

      {/* Scan button */}
      <div className="flex justify-center">
        <button
          onClick={handleScan}
          disabled={loading || !code.trim()}
          className="px-8 py-2.5 bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-semibold rounded-lg transition-colors text-sm shadow-sm"
        >
          {loading ? (
            <span className="flex items-center gap-2">
              <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
              </svg>
              Scanning…
            </span>
          ) : "Scan for References"}
        </button>
      </div>

      {/* Error */}
      {error && (
        <div className="rounded-lg border border-red-300 bg-red-50 px-4 py-3 text-red-700 text-sm font-mono">
          {error}
        </div>
      )}

      {/* Results */}
      {hits !== null && (
        hits.length === 0 ? (
          <div className="text-center py-10 text-green-600 font-medium text-sm">
            ✓ &nbsp;No references to breaking changes found in this file
          </div>
        ) : (
          <div className="space-y-3">
            {/* Summary */}
            <p className="text-sm text-gray-500">
              Found <span className="font-semibold text-red-700">{hits.length}</span> reference(s) across{" "}
              <span className="font-semibold">{distinctFiles}</span> file(s)
            </p>

            {/* Per-change groups */}
            {Object.entries(groups).map(([changePath, groupHits]) => (
              <div key={changePath} className="border border-red-200 rounded-lg overflow-hidden">
                <div className="bg-red-50 px-4 py-2 flex items-center gap-2">
                  <span className="text-red-600 text-sm">🔴</span>
                  <span className="text-sm font-semibold text-red-800">{changePath}</span>
                  <span className="ml-auto text-xs text-red-500">{groupHits.length} hit(s)</span>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="bg-gray-50 border-b border-gray-200">
                        <th className="px-4 py-2 text-left text-[11px] font-semibold uppercase tracking-wider text-gray-500 whitespace-nowrap">File</th>
                        <th className="px-4 py-2 text-left text-[11px] font-semibold uppercase tracking-wider text-gray-500">Line</th>
                        <th className="px-4 py-2 text-left text-[11px] font-semibold uppercase tracking-wider text-gray-500">Code</th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-100">
                      {groupHits.map((h, i) => (
                        <tr key={i} className="hover:bg-gray-50">
                          <td className="px-4 py-2.5 font-mono text-xs text-gray-500 whitespace-nowrap">{h.file}</td>
                          <td className="px-4 py-2.5 font-mono text-xs text-gray-400 whitespace-nowrap">{h.line_num}</td>
                          <td className="px-4 py-2.5 font-mono text-xs text-gray-800 whitespace-pre">{h.line}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            ))}
          </div>
        )
      )}
    </div>
  );
}
