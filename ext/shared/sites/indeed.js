window.Joblantern.registerSite({
  id: "indeed",
  name: "Indeed",
  match: (loc) =>
    /(^|\.)indeed\.[a-z.]+$/.test(loc.host) && loc.pathname.includes("viewjob"),
  extract: () => {
    const t = window.Joblantern.textOf;
    const title = t('[data-testid="jobsearch-JobInfoHeader-title"]') || t("h1.jobsearch-JobInfoHeader-title");
    const company = t('[data-testid="inlineHeader-companyName"]') || t(".jobsearch-CompanyInfoContainer a");
    const location = t('[data-testid="inlineHeader-companyLocation"]') || t(".jobsearch-JobInfoHeader-subtitle div");
    const desc = t("#jobDescriptionText");
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
    document.querySelector('[data-testid="jobsearch-JobInfoHeader-title"]') ||
    document.querySelector("h1.jobsearch-JobInfoHeader-title"),
});
