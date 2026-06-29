import { Post } from '@/types/post';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

const MOCK_POSTS: Post[] = [
  {
    id: '1',
    title: 'Tre giorni sulle Dolomiti',
    description: 'Un trekking tra le cime più iconiche delle Alpi orientali, con rifugi, marmotte e tramonti infuocati.',
    thumbnail: '',
    created_at: '2025-06-01T10:00:00Z',
  },
  {
    id: '2',
    title: 'Costarica in moto',
    description: 'Da San José alla penisola di Osa passando per vulcani attivi, foreste pluviali e spiagge deserte.',
    thumbnail: '',
    created_at: '2025-04-15T08:30:00Z',
  },
  {
    id: '3',
    title: 'Giappone fuori stagione',
    description: 'Tokyo, Kyoto e le Alpi giapponesi in febbraio: nessuna folla, neve fresca e onsen caldi.',
    thumbnail: '',
    created_at: '2025-02-20T14:00:00Z',
  },
];

export async function getPosts(): Promise<Post[]> {
  // TODO: sostituire con la chiamata reale quando l'endpoint /posts sarà pronto
  // const res = await fetch(`${BASE_URL}/posts`, { cache: 'no-store' });
  // if (!res.ok) throw new Error('Errore nel caricamento dei post');
  // return res.json();
  void BASE_URL;
  return MOCK_POSTS;
}
