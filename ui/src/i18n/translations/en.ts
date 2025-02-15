import { type Components, count, params } from "@nanostores/i18n"

const sharedErrorTranslations: Components = {
    Input: {
        "Invalid/Empty": params("{name} must not be empty"),
    },
    "ChangePassword/Errors": {
        CurrentPasswordIncorrect: "Incorrect password",
        EmptyCurrentPassword: "Please enter current password",
        EmptyNewPassword: "Please enter a new password",
        EmptyRepeateNewPassword: "Please repeat the new password",
        NewPasswordsDoNotMatch: "New passwords don't match",
        NewPasswordIsOldPassword: "New password can't be old password",
    },
}

export const translations = {
    "components/Navigation": {
        Settings: "Settings",
        Back: "Back",
    },

    "components/Memo": {
        ShowMore: "Show More",
    },

    "components/Memo/Actions": {
        View: "View",
        Edit: "Edit",
        Archive: "Archive",
        Delete: "Delete",
        Unarchive: "Unarchive",
        Restore: "Restore",
    },

    "components/Memo/DateTime": {
        CreatedAt: "Created at",
        UpdatedAt: "Updated at",
    },

    "components/MemoList/DayHeader": {
        Today: "Today",
        Yesterday: "Yesterday",
    },

    "components/MemoList/LayoutSelect": {
        Label: "Select list layout",
        LayoutMasonry: "Masonry",
        LayoutSingle: "Single",
        LayoutUltraCompact: "Ultra Compact",
    },

    "components/Editor": {
        Cancel: "Cancel",
        Save: "Save",
        DiscardChangesTitle: "Discard Changes",
        DiscardChangesDescription:
            "Are you sure you want to discard any changes?",
        DiscardChangesConfirmation: "Discard",
    },

    "components/Editor/Toolbar": {
        TextFormatting: "Text formatting",
        TextFormattingBold: "Bold",
        TextFormattingItalic: "Italic",
        TextFormattingMonospace: "Monospace",
        InsertLink: "Insert Link",
    },

    "components/DateTime": {
        datetime: params("{date} at {time}"),
        invalidTime: params(`Invalid date "{date}": {error}`),
    },

    "components/ThemeSwitcher": {
        ColoursDefault: "Default Colours",
        ColoursWarm: "Warm",
        RosePine: "Rosé Pine",
        Auto: "Auto",
        Light: "Light",
        Dark: "Dark",
        SelectColourSchemeAriaLabel: "Select theme",
        SelectModeAriaLabel: "Select light/dark mode",
    },

    "components/Notifications": {
        Label: "Notifications",
        Dismiss: "Dismiss",
    },

    "components/Greeting": {
        Morning: "Good Morning, ",
        Afternoon: "Good Afternoon, ",
        Evening: "Good Evening, ",
    },

    "components/MemoListHeader": {
        Archived: "Archived memos",
        Deleted: "Deleted Memos",
        MemosForTag: "Memos tagged",
        MemosForExactDate: "created on",
        MemosForExactDateStandalone: "Memos created on",
        MemosForQuery: "containing",
        MemosForQueryStandalone: "Memos containing",
    },

    "components/MemoListFilter": {
        TriggerLabel: "More filter",
        OffScreenTitle: "Tags",
        OffScreenDescription: "Tags and state filters",
    },

    "components/MemoListFilter/DatePicker": {
        Today: "Today",
    },

    "components/MemoListFilter/TagTree": {
        Label: "Tags",
    },

    "components/MemoListFilter/StateFilter": {
        Archived: "Archived",
        Deleted: "Deleted",
    },

    "pages/Errors/NotFound": {
        Title: "Not Found",
        Detail: "The requested page could not be found",
    },

    "pages/Errors/Unauthorized": {
        Title: "Unauthorized",
        Detail: "You are not authorized to see this page",
    },

    "pages/Login": {
        ...sharedErrorTranslations.Input,
        Title: "Login",
        UsernameLabel: "Username",
        PasswordLabel: "Password",
        LoginButton: "Login",
    },

    "pages/LoginChangePassword": {
        ...sharedErrorTranslations.Input,
        ...sharedErrorTranslations["ChangePassword/Errors"],
        Title: "Change Password",
        CurrentPasswordLabel: "Current Password",
        NewPasswordLabel: "New Password",
        RepeatNewPasswordLabel: "Repeat New Password",
        ChangePasswordButton: "Change Password",
    },

    "pages/MemoSinglePage": {
        Back: "Back",
    },

    "pages/Settings": {
        Title: "Settings",
    },

    "pages/Settings/InterfaceSettings": {
        Title: "Interface",
        Description: "Control how the user interface looks and behaves.",
        SectionTheme: "Theme",
        LabelColourScheme: "Colour Scheme",
        LabelModeOverride: "Light/Dark Mode Override",
        LabelIcon: "Icon",
        SectionControls: "Controls",
        LabelEnableVimKeybindings: "Enable Vim keybindings",
        LabelEnableDoubleClickToEdit: "Enable double click to edit memos",
    },

    "pages/Settings/LocaleSettings": {
        Title: "Language & Locale",
        Description: "Set your preferred language and locale.",
        LabelSelectLanguage: "Language",
        LabelSelectRegion: "Region",
    },

    "pages/Settings/AccountSettings": {
        ...sharedErrorTranslations.Input,
        ...sharedErrorTranslations["ChangePassword/Errors"],
        Title: "Account",
        Description: "Manage your account.",
        DisplayNameLabel: "Display Name",
        UpdateDisplayNameButton: "Update",
        ChangePasswordTitle: "Change Password",
        CurrentPasswordLabel: "Current Password",
        NewPasswordLabel: "New Password",
        RepeatNewPasswordLabel: "Repeat New Password",
        ChangePasswordButton: "Change",
        EmptyDisplayName: "Display Name must not be empty",
        Logout: "Logout",
    },

    "pages/Settings/SystemSettings": {
        Title: "System",
        SectionAPITokensTitle: "API Tokens",
        SectionAPITokensDescription: "Create and revoke API Access Tokens.",
    },

    "pages/Settings/SystemSettings/New": {
        Title: "Create API Token",
        FieldNameLabel: "Name",
        FieldExpiresInLabel: "Expires In",
        FieldExpiresInValueDays: count({
            one: "{count} day",
            many: "{count} days",
        }),
        FieldExpiresInValueMonths: count({
            one: "{count} month",
            many: "{count} months",
        }),
        ButtonLabel: "Create",
    },

    "pages/Settings/SystemSettings/LastCreated": {
        Title: "Created API Token",
        Notice: "Please note this token. IT WILL NOT BE SHOWN AGAIN!",
    },

    "pages/Settings/SystemSettings/List": {
        Title: "API Tokens",
        LabelName: "Name",
        LabelExpires: "Expires",
        LabelCreated: "Created",
        DeleteButton: "Delete",
        PrevPage: "Previous Page",
        NextPage: "Next Page",
    },
} satisfies Components
