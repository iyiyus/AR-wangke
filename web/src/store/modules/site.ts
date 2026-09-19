import { reactive } from 'vue'

export const siteConfig = reactive({
  sitename: '',
  logo: '',
  favicon: '',
  loaded: false
})

export async function loadSiteConfig() {
  if (siteConfig.loaded) return
  try {
    const r = await fetch('/api/site')
    const res = await r.json()
    const d = res?.data
    if (d) {
      if (d.sitename) {
        siteConfig.sitename = d.sitename
        document.title = d.sitename
      }
      if (d.logo) siteConfig.logo = d.logo
      if (d.favicon) {
        siteConfig.favicon = d.favicon
        let link = document.querySelector("link[rel='icon']") as HTMLLinkElement
        if (!link) {
          link = document.createElement('link')
          link.rel = 'icon'
          document.head.appendChild(link)
        }
        link.href = d.favicon
      }
    }
  } catch {}
  siteConfig.loaded = true
}
