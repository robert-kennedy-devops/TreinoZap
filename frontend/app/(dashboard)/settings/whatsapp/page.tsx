'use client';

import { useEffect, useRef, useState } from 'react';
import QRCode from 'qrcode';
import { api } from '@/lib/api';
import type { WhatsAppStatus } from '@/types';
import Header from '@/components/layout/Header';
import Button from '@/components/ui/Button';
import Badge from '@/components/ui/Badge';
import Card from '@/components/ui/Card';

export default function WhatsAppAdminPage() {
  const [status, setStatus] = useState<WhatsAppStatus | null>(null);
  const [qrData, setQrData] = useState('');
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [error, setError] = useState('');
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);

  async function loadStatus() {
    try {
      const s = await api.get<WhatsAppStatus>('/admin/whatsapp/status');
      setStatus(s);
      if (s.connected) {
        setQrData('');
        stopPolling();
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Acesso negado');
    } finally {
      setLoading(false);
    }
  }

  async function loadQR() {
    try {
      const res = await api.get<{ qr_code: string }>('/admin/whatsapp/qr');
      if (res.qr_code && res.qr_code !== qrData) {
        setQrData(res.qr_code);
      }
    } catch {
      // QR not ready yet or already connected
    }
  }

  function startPolling() {
    if (pollRef.current) return;
    pollRef.current = setInterval(async () => {
      await loadQR();
      await loadStatus();
    }, 3000);
  }

  function stopPolling() {
    if (pollRef.current) {
      clearInterval(pollRef.current);
      pollRef.current = null;
    }
  }

  // Render QR code onto canvas whenever qrData changes
  useEffect(() => {
    if (!qrData || !canvasRef.current) return;
    QRCode.toCanvas(canvasRef.current, qrData, {
      width: 256,
      margin: 2,
      color: { dark: '#000000', light: '#ffffff' },
    }).catch(() => {});
  }, [qrData]);

  useEffect(() => {
    loadStatus();
    return () => stopPolling();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleConnect() {
    setActionLoading(true);
    setError('');
    try {
      await api.post('/admin/whatsapp/connect');
      // Give the server a moment to generate the first QR
      setTimeout(() => {
        loadQR();
        loadStatus();
        startPolling();
      }, 1500);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao conectar');
    } finally {
      setActionLoading(false);
    }
  }

  async function handleDisconnect() {
    if (!confirm('Desconectar WhatsApp? Você poderá reconectar depois pelo painel.')) return;
    setActionLoading(true);
    stopPolling();
    try {
      await api.post('/admin/whatsapp/disconnect');
      setQrData('');
      await loadStatus();
    } finally {
      setActionLoading(false);
    }
  }

  if (loading) return <div className="p-6 text-sm text-gray-500">Carregando...</div>;

  return (
    <div>
      <Header title="WhatsApp — Configurações" />
      <div className="p-6 max-w-lg space-y-4">
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 text-sm px-4 py-3 rounded-md">
            {error}
          </div>
        )}

        <Card className="p-5 space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-sm font-semibold text-gray-800">Status da Conexão</h2>
            <Badge status={status?.connected ? 'connected' : 'disconnected'} />
          </div>

          {status?.phone && (
            <p className="text-sm text-gray-600">
              📱 +{status.phone}
            </p>
          )}
          {status?.jid && (
            <p className="text-xs text-gray-400 font-mono">JID: {status.jid}</p>
          )}
          {status?.last_connected && (
            <p className="text-xs text-gray-400">
              Última conexão: {new Date(status.last_connected).toLocaleString('pt-BR')}
            </p>
          )}

          <div className="flex gap-3 flex-wrap">
            {!status?.connected && (
              <Button onClick={handleConnect} loading={actionLoading}>
                Conectar
              </Button>
            )}
            {status?.connected && (
              <Button variant="danger" onClick={handleDisconnect} loading={actionLoading}>
                Desconectar
              </Button>
            )}
            <Button variant="secondary" onClick={loadStatus}>
              Atualizar
            </Button>
          </div>
        </Card>

        {qrData && (
          <Card className="p-5">
            <p className="text-sm font-semibold text-gray-800 mb-1">Escaneie com o WhatsApp</p>
            <p className="text-xs text-gray-500 mb-4">
              Abra o WhatsApp → Dispositivos vinculados → Vincular dispositivo → aponte para o QR abaixo
            </p>
            <div className="flex justify-center">
              <canvas
                ref={canvasRef}
                className="rounded-lg border border-gray-200"
              />
            </div>
            <p className="text-xs text-gray-400 text-center mt-3">
              O QR expira em ~20s e é renovado automaticamente
            </p>
          </Card>
        )}

        {!status?.connected && !qrData && (
          <Card className="p-5">
            <p className="text-sm text-gray-500">
              Clique em <strong>Conectar</strong> para gerar o QR Code de pareamento.
            </p>
          </Card>
        )}
      </div>
    </div>
  );
}
