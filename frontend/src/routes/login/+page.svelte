<script lang="ts">
    import { browser } from '$app/environment';
    import { toastStore } from '$lib/stores/toastStore';
  
    let isLogin = true;
    let username = '';
    let email = '';
    let password = '';
    let identifier = '';
    let error = '';
  
    const API_URL = browser 
      ? (window.location.hostname === "localhost"
        ? "http://localhost:8000"
        : "http://address_module_backend:8000")
      : "http://address_module_backend:8000";
  
    const toggleMode = () => {
      isLogin = !isLogin;
      error = '';
    };
  
    async function handleSubmit() {
      const url = isLogin ? '/auth/login' : '/auth/register';
      const payload = isLogin
        ? { identifier, password }
        : { username, email, password };
  
      try {
        const res = await fetch(`${API_URL}${url}`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
  
        if (!res.ok) {
          const errText = await res.text();
          toastStore.push(`Error: ${errText}`, 'error');
          throw new Error(errText);
        }
  
        const data = await res.json();
  
        if (isLogin && data.token) {
          localStorage.setItem('token', data.token);
          toastStore.push('Login successful!', 'success');
        } else {
          isLogin = true;
          toastStore.push('Registration successful! Please log in.', 'success');
        }
      } catch (err) {
        console.error(err);
        toastStore.push('Something went wrong. Please try again.', 'error');
      }
    }
  </script>
  
  <style>
    .auth-container {
      max-width: 360px;
      margin: 3rem auto;
      padding: 2rem;
      border-radius: 10px;
      box-shadow: 0 2px 10px rgba(0,0,0,0.1);
      background: #fff;
      font-family: sans-serif;
    }
  
    h2 {
      text-align: center;
      margin-bottom: 1.5rem;
    }
  
    input {
      width: 93%;
      padding: 0.75rem;
      margin-bottom: 1rem;
      border: 1px solid #ccc;
      border-radius: 5px;
    }
  
    button {
      width: 100%;
      padding: 0.75rem;
      background: #0077cc;
      color: #fff;
      border: none;
      border-radius: 5px;
      font-weight: bold;
      cursor: pointer;
    }
  
    .toggle {
      margin-top: 1rem;
      text-align: center;
      color: #0077cc;
      cursor: pointer;
      font-size: 0.9rem;
    }
  
    .error {
      color: red;
      margin-top: 1rem;
      text-align: center;
    }
  </style>
  
  <div class="auth-container">
    <h2>{isLogin ? 'Login' : 'Register'}</h2>
  
    {#if !isLogin}
      <input type="text" placeholder="Username" bind:value={username} />
      <input type="email" placeholder="Email" bind:value={email} />
    {:else}
      <input type="text" placeholder="Email or Username" bind:value={identifier} />
    {/if}
  
    <input type="password" placeholder="Password" bind:value={password} />
    <button on:click={handleSubmit}>{isLogin ? 'Login' : 'Register'}</button>
  
    {#if error}
      <div class="error">{error}</div>
    {/if}
  
    <div class="toggle" on:click={toggleMode}>
      {isLogin ? 'Don’t have an account? Register' : 'Already have an account? Login'}
    </div>
  </div>
  