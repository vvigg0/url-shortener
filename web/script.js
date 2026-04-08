const shortenForm = document.getElementById('shorten-form')
const urlInput = document.getElementById('url-input')
const customCodeInput = document.getElementById('custom-code-input')
const shortenResult = document.getElementById('shorten-result')

const redirectForm = document.getElementById('redirect-form')
const redirectLinkInput = document.getElementById('redirect-link')

const loadLinksBtn = document.getElementById('load-links-btn')
const allLinksResult = document.getElementById('all-links-result')

const analyticsForm = document.getElementById('analytics-form')
const analyticsCodeInput = document.getElementById('analytics-code')
const aggregateSelect = document.getElementById('aggregate-select')
const aggregateTrigger = document.getElementById('aggregate-trigger')
const aggregateTriggerText = document.getElementById('aggregate-trigger-text')
const aggregateMenu = document.getElementById('aggregate-menu')
const aggregateOptions = document.querySelectorAll('.custom-select-option')
const analyticsMeta = document.getElementById('analytics-meta')
const analyticsResult = document.getElementById('analytics-result')

shortenForm.addEventListener('submit', onShorten)
redirectForm.addEventListener('submit', onRedirect)
loadLinksBtn.addEventListener('click', onLoadAllLinks)
analyticsForm.addEventListener('submit', onAnalytics)

let aggregateModeValue = ''

aggregateTrigger.addEventListener('click', toggleAggregateMenu)

aggregateOptions.forEach(option => {
  option.addEventListener('click', () => {
    aggregateModeValue = option.dataset.value || ''
    aggregateTriggerText.textContent = option.textContent
    aggregateMenu.classList.add('hidden')

    aggregateOptions.forEach(item => item.classList.remove('active'))
    option.classList.add('active')
  })
})

document.addEventListener('click', event => {
  if (!aggregateSelect.contains(event.target)) {
    aggregateMenu.classList.add('hidden')
  }
})

function toggleAggregateMenu() {
  aggregateMenu.classList.toggle('hidden')
}



function normalizeCode(value) {
  return value.trim()
}

function createMessage(text, type = 'success') {
  return `<div class="message ${type}">${escapeHtml(text)}</div>`
}

function extractError(payload, fallback) {
  if (!payload) return fallback
  if (typeof payload.err === 'string') return payload.err
  if (typeof payload.error === 'string') return payload.error
  return fallback
}

function buildFullShortLink(payload) {
  const raw = payload?.res

  if (typeof raw === 'string' && raw.trim()) {
    return raw.trim()
  }

  if (payload?.short_code) {
    return `localhost:8080/s/${payload.short_code}`
  }

  throw new Error('Сервер вернул некорректный ответ')
}

async function onShorten(event) {
  event.preventDefault()
  shortenResult.classList.remove('hidden')
  shortenResult.innerHTML = createMessage('Создаю короткую ссылку...')

  try {
    const response = await fetch('/shorten', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        url: urlInput.value.trim(),
        custom_code: customCodeInput.value.trim()
      })
    })

    const payload = await response.json().catch(() => null)

    if (!response.ok) {
      throw new Error(extractError(payload, 'Не удалось сократить ссылку'))
    }

    const shortLink = buildFullShortLink(payload)
    const shortCode = getShortCodeFromLink(shortLink)
    const clickableLink = shortLink.startsWith('http://') || shortLink.startsWith('https://')
      ? shortLink
      : `http://${shortLink}`

    shortenResult.innerHTML = `
      <p class="result-title">Короткая ссылка успешно создана</p>
      <div class="result-link">
        <a href="${escapeAttribute(clickableLink)}" target="_blank" rel="noopener noreferrer">${escapeHtml(shortLink)}</a>
        <div class="inline-actions">
          <button class="mini-btn" type="button" data-copy="${escapeAttribute(shortLink)}">Скопировать</button>
          <button class="mini-btn" type="button" data-fill-code="${escapeAttribute(shortCode)}">В аналитику</button>
        </div>
      </div>
    `

    shortenResult.querySelector('[data-copy]')?.addEventListener('click', async e => {
      const value = e.currentTarget.dataset.copy || ''
      await navigator.clipboard.writeText(value)
      e.currentTarget.textContent = 'Скопировано'
    })

    shortenResult.querySelector('[data-fill-code]')?.addEventListener('click', e => {
      const code = e.currentTarget.dataset.fillCode || ''
      analyticsCodeInput.value = code
      redirectLinkInput.value = shortLink
      analyticsCodeInput.focus()
    })
  } catch (error) {
    shortenResult.innerHTML = createMessage(error.message || 'Ошибка при создании ссылки', 'error')
  }
}

function onRedirect(event) {
  event.preventDefault()

  let link = redirectLinkInput.value.trim()
  if (!link) return

  if (!link.startsWith('http://') && !link.startsWith('https://')) {
    link = `http://${link}`
  }

  try {
    const url = new URL(link)
    window.open(url.toString(), '_blank', 'noopener,noreferrer')
  } catch {
    alert('Введи корректную ссылку')
  }
}

async function onLoadAllLinks() {
  allLinksResult.classList.remove('hidden')
  allLinksResult.innerHTML = createMessage('Загружаю ссылки...')

  try {
    const response = await fetch('/shortened')
    const payload = await response.json().catch(() => null)

    if (!response.ok) {
      throw new Error(extractError(payload, 'Не удалось получить список ссылок'))
    }

    const links = payload?.res?.links

    if (!Array.isArray(links)) {
      throw new Error('Сервер вернул некорректный список ссылок')
    }

    if (!links.length) {
      allLinksResult.innerHTML = createMessage('Список ссылок пуст', 'error')
      return
    }

    allLinksResult.innerHTML = `
      <table class="table">
        <thead>
          <tr>
            <th>#</th>
            <th>Короткая ссылка</th>
            <th>Полная ссылка</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          ${links.map((link, index) => {
            const shortLink = `localhost:8080/s/${link.short_code}`
            const clickableLink = `http://${shortLink}`

            return `
              <tr>
                <td>${escapeHtml(String(index + 1))}</td>
                <td>
                  <a href="${escapeAttribute(clickableLink)}" target="_blank" rel="noopener noreferrer">
                    ${escapeHtml(shortLink)}
                  </a>
                </td>
                <td>${escapeHtml(link.full_url || '—')}</td>
                <td>
                  <div class="inline-actions">
                    <button
                      class="mini-btn"
                      type="button"
                      data-copy-link="${escapeAttribute(shortLink)}"
                    >
                      Скопировать
                    </button>
                    <button
                      class="mini-btn"
                      type="button"
                      data-fill-analytics="${escapeAttribute(link.short_code || '')}"
                    >
                      В аналитику
                    </button>
                  </div>
                </td>
              </tr>
            `
          }).join('')}
        </tbody>
      </table>
    `

    allLinksResult.querySelectorAll('[data-copy-link]').forEach(button => {
      button.addEventListener('click', async e => {
        const value = e.currentTarget.dataset.copyLink || ''
        await navigator.clipboard.writeText(value)
        e.currentTarget.textContent = 'Скопировано'
      })
    })

    allLinksResult.querySelectorAll('[data-fill-analytics]').forEach(button => {
      button.addEventListener('click', e => {
        const code = e.currentTarget.dataset.fillAnalytics || ''
        analyticsCodeInput.value = code
        analyticsCodeInput.focus()
      })
    })
  } catch (error) {
    allLinksResult.innerHTML = createMessage(error.message || 'Ошибка при загрузке списка ссылок', 'error')
  }
}

async function onAnalytics(event) {
  event.preventDefault()

  const code = normalizeCode(analyticsCodeInput.value)
  const mode = aggregateModeValue

  analyticsMeta.classList.remove('hidden')
  analyticsResult.classList.remove('hidden')
  analyticsMeta.innerHTML = createMessage('Загружаю аналитику...')
  analyticsResult.innerHTML = ''

  try {
    const query = mode ? `?aggregate_by=${encodeURIComponent(mode)}` : ''
    const response = await fetch(`/analytics/${encodeURIComponent(code)}${query}`)
    const payload = await response.json().catch(() => null)

    if (!response.ok) {
      throw new Error(extractError(payload, 'Не удалось получить аналитику'))
    }

    const data = payload?.res
    renderAnalyticsMeta(data)
    renderAnalyticsTable(data, mode)
  } catch (error) {
    analyticsMeta.innerHTML = createMessage(error.message || 'Ошибка при загрузке аналитики', 'error')
    analyticsResult.innerHTML = ''
  }
}

function renderAnalyticsMeta(data) {
  if (!data) {
    analyticsMeta.innerHTML = createMessage('Пустой ответ от сервера', 'error')
    return
  }

  const blocks = [
    { label: 'Short code', value: data.short_code || '—' },
    { label: 'Полная ссылка', value: data.full_url || '—' }
  ]

  if (typeof data.clicks === 'number') {
    blocks.push({ label: 'Всего кликов', value: String(data.clicks) })
  } else if (Array.isArray(data.data)) {
    const total = data.data.reduce((sum, item) => sum + Number(item.clicks || 0), 0)
    blocks.push({ label: 'Суммарно по выборке', value: String(total) })
  }

  analyticsMeta.innerHTML = `<div class="meta-grid">${blocks.map(item => `
    <div class="meta-box">
      <span>${escapeHtml(item.label)}</span>
      <strong>${escapeHtml(item.value)}</strong>
    </div>
  `).join('')}</div>`
}

function renderAnalyticsTable(data, mode) {
  if (!data) {
    analyticsResult.innerHTML = ''
    return
  }

  if (Array.isArray(data.visits)) {
    if (!data.visits.length) {
      analyticsResult.innerHTML = createMessage('Переходов пока нет')
      return
    }

    analyticsResult.innerHTML = buildTable(
      ['#', 'Visited at', 'User-Agent'],
      data.visits.map((visit, index) => [
        String(index + 1),
        formatDateTime(visit.visited_at),
        visit['user-agent'] || visit.user_agent || '—'
      ])
    )
    return
  }

  if (Array.isArray(data.data) && mode === 'user-agent') {
    if (!data.data.length) {
      analyticsResult.innerHTML = createMessage('Переходов пока нет')
      return
    }

    analyticsResult.innerHTML = buildTable(
      ['#', 'User-Agent', 'Clicks'],
      data.data.map((item, index) => [
        String(index + 1),
        item.user_agent || '—',
        String(item.clicks || 0)
      ])
    )
    return
  }

  if (Array.isArray(data.data)) {
    if (!data.data.length) {
      analyticsResult.innerHTML = createMessage('Переходов пока нет')
      return
    }

    analyticsResult.innerHTML = buildTable(
      ['#', mode === 'month' ? 'Month' : 'Day', 'Clicks'],
      data.data.map((item, index) => [
        String(index + 1),
        item.period || '—',
        String(item.clicks || 0)
      ])
    )
    return
  }

  if (typeof data.clicks === 'number' && data.clicks === 0) {
    if (mode === 'user-agent') {
      analyticsResult.innerHTML = buildTable(
        ['#', 'User-Agent', 'Clicks'],
        [['1', '—', '0']]
      )
      return
    }

    if (mode === 'month') {
      analyticsResult.innerHTML = buildTable(
        ['#', 'Month', 'Clicks'],
        [['1', '—', '0']]
      )
      return
    }

    if (mode === 'day') {
      analyticsResult.innerHTML = buildTable(
        ['#', 'Day', 'Clicks'],
        [['1', '—', '0']]
      )
      return
    }

    analyticsResult.innerHTML = createMessage('Переходов пока нет')
    return
  }

  analyticsResult.innerHTML = createMessage('Переходов пока нет')
}

function buildTable(headers, rows) {
  if (!rows.length) {
    return createMessage('По этой ссылке пока нет данных', 'error')
  }

  return `
    <table class="table">
      <thead>
        <tr>${headers.map(header => `<th>${escapeHtml(header)}</th>`).join('')}</tr>
      </thead>
      <tbody>
        ${rows.map(row => `<tr>${row.map(cell => `<td>${escapeHtml(cell)}</td>`).join('')}</tr>`).join('')}
      </tbody>
    </table>
  `
}

function getShortCodeFromLink(link) {
  try {
    const url = new URL(link)
    return url.pathname.split('/').filter(Boolean).pop() || ''
  } catch {
    return ''
  }
}

function formatDateTime(value) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  }).format(date)
}

function escapeHtml(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;')
}

function escapeAttribute(value) {
  return escapeHtml(value)
}