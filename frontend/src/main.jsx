import React, { useState } from 'react';
import { createRoot } from 'react-dom/client';

function App() {
  const [letter, setLetter] = useState('');
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function submit(event) {
    event.preventDefault();
    setResult(null);
    setError('');
    if (!/^\p{L}$/u.test(letter)) {
      setError('Enter exactly one letter, such as L, Ł, or É.');
      return;
    }
    setLoading(true);
    try {
      const response = await fetch(`/cities/count?first_letter=${encodeURIComponent(letter)}`);
      if (!response.ok) throw new Error((await response.text()).trim() || 'Unable to count cities. Please try again.');
      const data = await response.json();
      setResult({ count: data.count, letter });
    } catch (err) {
      setError(err instanceof TypeError ? 'Could not connect. Please try again.' : err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main>
      <form onSubmit={submit} aria-busy={loading}>
        <label htmlFor="letter">First letter: </label>
        <input id="letter" value={letter} required disabled={loading}
          autoComplete="off" spellCheck={false}
          onChange={(event) => { setLetter(event.target.value); setResult(null); setError(''); }} />
        {' '}
        <button disabled={loading}>{loading ? 'Counting…' : 'Count cities'}</button>
      </form>
      <div role="status" aria-live="polite">
        {error && <p>{error}</p>}
        {result && <p>Cities starting with “{result.letter}”: {result.count}</p>}
      </div>
    </main>
  );
}

createRoot(document.getElementById('root')).render(<App />);
