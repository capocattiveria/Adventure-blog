'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function DashboardPage() {
  const [email, setEmail] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
      return;
    }
    // decodifica il payload del JWT (base64) per leggere l'email
    const payload = JSON.parse(atob(token.split('.')[1]));
    setEmail(payload.email);
  }, [router]);

  function handleLogout() {
    localStorage.removeItem('token');
    router.push('/login');
  }

  return (
    <main className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-2xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-2xl font-bold text-gray-900">Adventure Blog</h1>
          <button
            onClick={handleLogout}
            className="text-sm text-gray-500 hover:text-red-500 transition-colors"
          >
            Esci
          </button>
        </div>
        <p className="text-gray-700">Benvenuto, <span className="font-semibold">{email}</span></p>
      </div>
    </main>
  );
}
