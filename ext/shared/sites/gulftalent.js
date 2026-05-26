window.Joblantern.registerSite({
  id: "gulftalent",
  name: "GulfTalent",
  match: (loc) =>
    loc.host === "www.gulftalent.com" && /\/[a-z]{2}\/jobs\//.test(loc.pathname),
  extract: () => {
    const t = window.Joblantern.textOf;
    const title = t("h1.job-title") || t("h1#job-title");
    const company = t(".company-link") || t(".company-name a");
    const location = t(".job-location") || t(".job-meta .location");
    const desc = t(".job-description") || t("#job-description");
    if (!title && !company) return null;
    return {
      listing_url: window.location.href,
      listing_text: desc,
      company_name: company,
      claimed_address: location,
      role: title,
    };
  },
  anchorSelector: () =>
    document.querySelector("h1.job-title") || document.querySelector("h1#job-title"),
});
