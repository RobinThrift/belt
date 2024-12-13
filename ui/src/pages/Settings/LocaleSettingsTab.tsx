import { Select } from "@/components/Select"
import {
    type Language,
    type Region,
    supportedLanguages,
    supportedRegions,
    useT,
} from "@/i18n"
import { useSetting } from "@/state/settings"
import React, { useCallback, useMemo } from "react"

export const LocaleSettingsTab = React.forwardRef<HTMLDivElement>(
    function LocaleSettingsTab(_, forwardedRef) {
        let t = useT("pages/Settings/LocaleSettingsTab")

        let [lang, setLang] = useSetting("locale.language")
        let [region, setRegion] = useSetting("locale.region")

        let onChangeLanguage = useCallback(
            (v: Language) => {
                setLang(v)
            },
            [setLang],
        )

        let onChangeRegion = useCallback(
            (v: Region) => {
                setRegion(v)
            },
            [setRegion],
        )

        let languages = useMemo(
            () =>
                supportedLanguages.map((lang) => (
                    <Select.Option key={lang} value={lang}>
                        {new Intl.DisplayNames([lang], {
                            type: "language",
                        }).of(lang.toUpperCase())}
                    </Select.Option>
                )),
            [],
        )

        let regions = useMemo(
            () =>
                supportedRegions.map((locale) => (
                    <Select.Option key={locale} value={locale}>
                        {new Intl.DisplayNames([lang], {
                            type: "region",
                        }).of(locale.toUpperCase())}
                    </Select.Option>
                )),
            [lang],
        )

        return (
            <div ref={forwardedRef} className="settings-tab">
                <div className="settings-tab-section">
                    <h2>{t.Title}</h2>
                    <small>{t.Description}</small>
                </div>

                <div className="settings-tab-section">
                    <div className="sm:mb-0 md:grid grid-cols-6 space-y-1">
                        <label
                            htmlFor="language"
                            className="flex items-center font-semibold text-sm"
                        >
                            {t.LabelSelectLanguage}
                        </label>
                        <Select
                            className="col-span-5"
                            name="language"
                            ariaLabel={t.LabelSelectLanguage}
                            onChange={onChangeLanguage}
                            value={lang}
                        >
                            {languages}
                        </Select>
                    </div>

                    <div className="sm:mb-0 md:grid grid-cols-6 space-y-1">
                        <label
                            htmlFor="region"
                            className="flex items-center font-semibold text-sm"
                        >
                            {t.LabelSelectRegion}
                        </label>
                        <Select
                            className="col-span-5"
                            name="region"
                            ariaLabel={t.LabelSelectRegion}
                            onChange={onChangeRegion}
                            value={region}
                        >
                            {regions}
                        </Select>
                    </div>
                </div>
            </div>
        )
    },
)
