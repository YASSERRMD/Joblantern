// Package workbench renders the investigator workbench UI used by
// NGOs and journalists. The UI is intentionally minimal — the value
// is the data, not the chrome.
package workbench

const Page = `<!doctype html>
<html><head><meta charset="utf-8"><title>Joblantern — Investigator Workbench</title>
<style>
body{font-family:system-ui,sans-serif;margin:0;display:grid;grid-template-columns:280px 1fr;height:100vh}
aside{background:#0e1a2b;color:#e0eaff;padding:1rem;overflow:auto}
main{padding:1rem}
ul{list-style:none;padding:0}
li{padding:.4rem 0;border-bottom:1px solid #1c3046}
button{background:#3b6;color:#000;border:0;padding:.4rem .8rem;cursor:pointer;border-radius:.3rem}
</style></head><body>
<aside>
  <h2>Components</h2>
  <ul id="components"></ul>
  <button onclick="refresh()">Refresh</button>
</aside>
<main>
  <h1>Recruitment network</h1>
  <div id="graph" style="width:100%;height:80vh;border:1px solid #ddd"></div>
</main>
<script>
async function refresh(){
  const r = await fetch('/api/v1/graph/components').then(x=>x.json());
  const list = document.getElementById('components');
  list.innerHTML = '';
  r.forEach(c=>{
    const li = document.createElement('li');
    li.textContent = c.id + ' (' + c.size + ' nodes)';
    li.onclick = ()=>load(c.id);
    list.appendChild(li);
  });
}
function load(id){
  document.getElementById('graph').textContent = 'Component ' + id + ' selected. d3 layout served from /api/v1/graph/' + id;
}
refresh();
</script></body></html>`
