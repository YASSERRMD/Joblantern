window.Joblantern.registerSite({
  id: "linkedin",
  name: "LinkedIn",
  match: (loc) =>
    loc.host === "www.linkedin.com" && loc.pathname.includes("/jobs/"),
  extract: () => {
    const t = window.Joblantern.textOf;
    const title = t("h1.top-card-layout__title") || t("h1.job-details-jobs-unified-top-card__job-title");
    const company =
      t("a.topcard__org-name-link") ||
      t(".job-details-jobs-unified-top-card__company-name a") ||
      t(".job-details-jobs-unified-top-card__company-name");
    const location =
      t(".topcard__flavor.topcard__flavor--bullet") ||
      t(".job-details-jobs-unified-top-card__bullet");
    if (!title && !company) return null;
    return {
      listing_url: window.location.href,
      listing_text: t(".description__text") || t("#job-details") || "",
      company_name: company,
      claimed_address: location,
      role: title,
    };
  },
  anchorSelector: () =>
    document.querySelector(".jobs-unified-top-card__content--two-pane") ||
    document.querySelector("h1.top-card-layout__title") ||
    document.querySelector("h1.job-details-jobs-unified-top-card__job-title"),
});
