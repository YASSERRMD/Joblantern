window.Joblantern.registerSite({
  id: "naukrigulf",
  name: "Naukrigulf",
  match: (loc) =>
    loc.host === "www.naukrigulf.com" && loc.pathname.startsWith("/job-listing-"),
  extract: () => {
    const t = window.Joblantern.textOf;
    const title = t("h1.title") || t("h1.jdHeader");
    const company = t("a.companyName") || t(".company-name");
    const location = t(".jdCol-location") || t(".location");
    const desc = t(".job-description") || t(".jdDesc");
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
    document.querySelector("h1.title") || document.querySelector("h1.jdHeader"),
});
