window.Joblantern.registerSite({
  id: "jobstreet",
  name: "Jobstreet",
  match: (loc) =>
    /(^|\.)jobstreet\.[a-z.]+$/.test(loc.host) && loc.pathname.includes("/job/"),
  extract: () => {
    const t = window.Joblantern.textOf;
    const title = t('[data-automation="job-detail-title"]') || t("h1");
    const company = t('[data-automation="advertiser-name"]') || t('[data-automation="company-name"]');
    const location = t('[data-automation="job-detail-location"]');
    const desc = t('[data-automation="jobAdDetails"]');
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
    document.querySelector('[data-automation="job-detail-title"]') ||
    document.querySelector("h1"),
});
