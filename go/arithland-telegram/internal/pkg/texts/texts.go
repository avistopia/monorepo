package texts

var (
	Common_Next = "⬅️"
	Common_Back = "➡️"

	Common_Cancel    = "لغو ↪️"
	Common_Cancelled = "لغو شد"
)

var (
	Err_Internal          = "خطایی پیش آمد. تیم فنی شهر ریاضی در حال تلاش برای حل خطاست."
	Err_Forbidden         = "دسترسی به این قابلیت برای شما امکان‌پذیر نیست."
	Err_InvalidInteger    = "لطفا یک عدد وارد کنید."
	Err_InsufficientFunds = "موجودی کافی نیست."
	Err_NoQuestionsFound  = "پرسشی یافت نشد."
)

var (
	MainMenu_WelcomeToArithland = "به ربات شهر ریاضی خوش آمدید"
	MainMenu_DefaultResponse    = "درود! این ربات شهر ریاضی است. لطفا عملیات مورد نظر خود را از منوی کیبورد انتخاب کنید."
)

var (
	ArithlandConstitution_Show = "نمایش قانون اساسی شهر ریاضی 📜"

	ArithlandConstitution = "متن قانون شهر ریاضی در اینجا"
)

var (
	ProfileManagement_Show = "مدیریت پروفایل 👤"

	ProfileManagement = `
نام کاربری شما: {displayName}
موجودی شما: {balance} پروف
	`
	ProfileManagement_ChangeDisplayName     = "تغییر نام کاربری 📝"
	ProfileManagement_EnterDisplayName      = "لطفا نام کاربری جدید را ارسال کنید."
	ProfileManagement_DisplayNameIsTooShort = "نام کاربری باید دست کم ۲ کاراکتر باشد."
	ProfileManagement_DisplayNameIsTooLong  = "نام کاربری نباید بیش از ۳۰ کاراکتر باشد."
)

var (
	QuestionsManagement_Show = "مدیریت پرسش‌ها 💡"

	QuestionsManagement = `
شما تا کنون {questionsCount} پرسش خریده‌اید و به {answeredQuestionsCount} پرسش، پاسخ داده‌اید.
موجودی شما: {balance} پروف
	`

	QuestionsManagement_ShowUnansweredQuestions = "📝 نمایش پرسش‌های بی‌پاسخ من"
	QuestionsManagement_ShowAnsweredQuestions   = "📄 نمایش پرسش‌های پاسخ‌داده‌شده من"
	QuestionsManagement_BuyQuestion             = "{levelEmoji} خرید پرسش {level} ({price} پروف)"
	QuestionsManagement_Details                 = `
{header}

#️⃣ شماره پرسش: {id}

{levelEmoji} سطح: {level}

📜 پرسش:
{text}

✍️ پاسخ:
{answer}

زمان خرید سؤال:
{createdAt}

زمان پاسخ به سؤال:
{answeredAt}

{footer}
`
	QuestionsManagement_DetailsHeaderUnanswered = "📝 پرسش‌های بی‌پاسخ"
	QuestionsManagement_DetailsHeaderAnswered   = "📄 پرسش‌های پاسخ‌داده‌شده"
	QuestionsManagement_Unanswered              = "هنوز پاسخ داده نشده"

	QuestionsManagement_DetailsFooterUnanswered = "می‌توانید پاسخ این پرسش را ارسال کنید."
	QuestionsManagement_AlreadyAnswered         = "این سؤال قبلا پاسخ داده شده است."
	QuestionsManagement_IncorrectAnswer         = "متأسفانه پاسخ واردشده صحیح نیست. ❌"
	QuestionsManagement_CorrectAnswer           = `
پاسخ واردشده درست است. ✅

{price} پروف به حساب شما واریز شد.
`
)

var (
	AdminQuestions_Show = "ادمین پرسش‌ها 🏆"

	AdminQuestions = `
تعداد کل پرسش‌ها: {questionsCount}

برای مدیریت یک پرسش،لطفا شماره آن را وارد کنید.
`
	AdminQuestions_Create   = "{levelEmoji} افزودن پرسش {level}"
	AdminQuestions_NotFound = "پرسش یافت نشد"
	AdminQuestions_Details  = `
فعال : {active}

#️⃣ شماره پرسش: {id}

{levelEmoji} سطح: {level}

📜 پرسش:
{text}

✍️ پاسخ:
{answer}
`
	AdminQuestions_Created               = "یک پرسش {level} جدید با شماره {id} ایجاد شد. لطفا {field} پرسش را بفرستید."
	AdminQuestions_SendField             = "لطفا {field} پرسش را بفرستید."
	AdminQuestions_FieldUpdated          = "{field} پرسش ذخیره شد."
	AdminQuestions_FieldUpdatedSendField = "{updatedField} پرسش ذخیره شد. لطفا {nextField} پرسش را بفرستید."
	AdminQuestions_UpdateField           = "✏️ تغییر {field} پرسش"
	AdminQuestions_Delete                = "❌ حذف پرسش"
	AdminQuestions_UndoDelete            = "↩️ بازگردانی پرسش"
	AdminQuestions_Deleted               = "حذف شد"
	AdminQuestions_UndoDeleted           = "بازگردانی شد"
	AdminQuestions_Refresh               = "🔄 ریفرش"
	AdminQuestions_UpdateLevel           = "{levelEmoji} تغییر سطح پرسش به {level}"
	AdminQuestions_Level                 = "سطح"
)
