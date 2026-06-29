'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { login } from '@/services/authService';

export default function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  async function handleSubmit(e: { preventDefault(): void }) {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const { token } = await login({ email, password });
      localStorage.setItem('token', token);
      router.push('/dashboard');
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-sm">
      <h1 className="text-2xl font-bold text-center text-gray-900">Adventure Blog</h1>

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

      {error && <p className="text-red-500 text-sm text-center">{error}</p>}

      <button
        type="submit"
        disabled={loading}
        className="bg-blue-600 text-white rounded-lg py-2 font-semibold hover:bg-blue-700 disabled:opacity-60 transition-colors"
      >
        {loading ? 'Accesso in corso...' : 'Accedi'}
      </button>

      <p className="text-sm text-center text-gray-500">
        Non hai un account?{' '}
        <a href="/register" className="text-blue-600 underline">
          Registrati
        </a>
      </p>
    </form>
  );
}
