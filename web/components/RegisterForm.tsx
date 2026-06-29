'use client';

import { useState } from 'react';
import { register } from '@/services/authService';

export default function RegisterForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  async function handleSubmit(e: { preventDefault(): void }) {
    e.preventDefault();
    if (password !== confirm) {
      setError('Le password non coincidono');
      return;
    }
    setLoading(true);
    setError(null);

    try {
      await register({ email, password });
      setSuccess(true);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  if (success) {
    return (
      <div className="flex flex-col items-center gap-4 w-full max-w-sm">
        <p className="text-green-600 font-semibold">Account creato con successo!</p>
        <a href="/login" className="text-blue-600 underline text-sm">
          Vai al login
        </a>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-sm">
      <h1 className="text-2xl font-bold text-center text-gray-900">Crea un account</h1>

      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
        className="border border-gray-300 rounded-lg px-4 py-2 text-sm text-gray-900 bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />

      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
        className="border border-gray-300 rounded-lg px-4 py-2 text-sm text-gray-900 bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />

      <input
        type="password"
        placeholder="Conferma password"
        value={confirm}
        onChange={(e) => setConfirm(e.target.value)}
        required
        className="border border-gray-300 rounded-lg px-4 py-2 text-sm text-gray-900 bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />

      {error && <p className="text-red-500 text-sm text-center">{error}</p>}

      <button
        type="submit"
        disabled={loading}
        className="bg-blue-600 text-white rounded-lg py-2 font-semibold hover:bg-blue-700 disabled:opacity-60 transition-colors"
      >
        {loading ? 'Registrazione in corso...' : 'Registrati'}
      </button>

      <p className="text-sm text-center text-gray-500">
        Hai già un account?{' '}
        <a href="/login" className="text-blue-600 underline">
          Accedi
        </a>
      </p>
    </form>
  );
}
