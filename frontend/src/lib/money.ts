export function parseMoney(value: string | number): number {
    if (typeof value === "number") {
        return Math.floor(value); // Assume it's already cents if number, or just pass through. 
        // But context implies input is usually string "12.34". 
        // If the user types "12", they mean $12.00, so 1200 cents.
        // If raw number comes in as 12.34 (float), we still have the float risk.
        // Let's coerce to string first to be safe if it looks like a float.
        return parseMoney(String(value));
    }

    if (!value) return 0;

    let s = value.trim();
    const sign = s.startsWith("-") ? -1 : 1;
    if (s.startsWith("-") || s.startsWith("+")) {
        s = s.substring(1);
    }

    // Split on decimal point
    const parts = s.split(".");

    let whole = parts[0] || "0";
    let fraction = parts[1] || "00";

    // Pad or truncate fraction to 2 digits
    if (fraction.length === 0) fraction = "00";
    else if (fraction.length === 1) fraction = fraction + "0";
    else if (fraction.length > 2) fraction = fraction.substring(0, 2);

    // Combine and parse integer
    const combined = whole + fraction;

    // Remove leading zeros to avoid octal confusion (though parseInt shouldn't if no 0x)
    // but primarily to handle "0.50" -> "0" + "50" -> "050" -> 50 correctly.
    const cents = parseInt(combined, 10);

    return isNaN(cents) ? 0 : sign * cents;
}
