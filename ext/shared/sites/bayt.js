window.Joblantern.registerSite({
  id: "bayt",
  name: "Bayt",
  match: (loc) =>
    loc.host === "www.bayt.com" && /\/(en|ar)\/.*\/jobs\//.test(loc.pathname),
  extract: () => {
    const t = window.Joblantern.textOf;
    const title = t("h1.job-title") || t("h1#job_title");
    const company =
      t('a[href*="/companies/"]') || t(".company-name a") || t(".t-company");
    const location = t(".t-mute a[href*='/locations/']") || t(".t-location");
    const desc = t("#job_description") || t(".job-description");
    if (!title && !company) return null;
    return {
      listing_url: window.location.href,
      listing_text: desc,
      company_name: company,
      claimed_address: location,
      role: title,
      jurisdiction: "AE",
    };
  },
  anchorSelector: () =>
    document.querySelector("h1.job-title") || document.querySelector("h1#job_title"),
});
