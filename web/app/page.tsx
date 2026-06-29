import PostCard from '@/components/PostCard';
import { getPosts } from '@/services/postService';

export default async function HomePage() {
  const posts = await getPosts();

  return (
    <main className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-100">
        <div className="max-w-5xl mx-auto px-6 py-5 flex items-center justify-between">
          <h1 className="text-xl font-bold text-gray-900 tracking-tight">Adventure Blog</h1>
          <nav className="flex items-center gap-4 text-sm">
            <a href="/login" className="text-gray-500 hover:text-gray-900 transition-colors">Accedi</a>
          </nav>
        </div>
      </header>

      <section className="max-w-5xl mx-auto px-6 py-12">
        {posts.length === 0 ? (
          <p className="text-center text-gray-400 py-24">Nessun post ancora.</p>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}
      </section>
    </main>
  );
}
