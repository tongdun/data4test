### Common Render Variables:
    Note: Support for automatic data generation: Common render variables (can refer to variables supported by Jmeter).
          Variable composition supports English, numbers, underscores, hyphens, e.g.: {name_test-1}


#### Context Variable References:
- '*{Name}*': Refers to the Name variable from the previous context, treated as a whole in JSON format.
- '**{Name}**': Refers to the Name variable from the previous context, treated as a whole in string format.
- '*{Name[2:3]}*': (To be incorporated) Refers to the values from the 3rd to the 4th position (inclusive of the 3rd, exclusive of the 4th due to Python slicing convention) of the Name variable from the previous context, treated as a whole in JSON format. Note: This exact syntax may need adjustment based on the programming language or framework being used, as Python slicing is 0-indexed and does not include the end index.
- '**{Name[2:3]}**': (To be incorporated) Refers to the values from the 3rd to the 4th position (similar caveats as above) of the Name variable from the previous context, treated as a whole in string format. This again may require adjustment based on the specific context and language being used.
- XXName: '{self}': Refers to the value of the XXName variable from the previous context. **Note**: The use of '{self}' as a placeholder for the variable's own value is being deprecated. It is recommended to replace all existing instances of '{self}' with the specific variable name as soon as possible.

#### Product Environment-Specific Parameters:
- Define some special parameters in the product configuration, in JSON format, e.g.: {"name1": "value1", "name2": "value2"}

#### System Built-In Parameters Supporting Multiple Languages:
- {Name}
- {Email}
- {Phone}
- {Mobile}
- {Contact}
- {Gender}
- {Sex}
- {CardNo}
- {BankNo}
- {SSN}
- {IDNo}
- {Address}
- {Company}
- {Country}
- {Province}
- {State}
- {City}

#### Range Data Auto-generation:
- {Int(a,b)}: Numbers, e.g.: {Int(1000,2000)}, generates a random number between 1000 and 2000.
- {Date(a,b)}: Date, e.g.: {Date(-30,1)}, generates a random date from 30 days ago to 1 day ahead.
- {Time(a,b)}: Date with time, e.g.: {Time(-30,1)}, generates a random date with time from 30 days ago to 1 day ahead.

#### Specified Length Data Auto-generation:
- {Str(length)}: Strings, e.g.: {Str(64)}, generates a string of length 64.
- {Rune(length)}: Chinese characters, e.g.: {Rune(128)}, generates a string of 128 Chinese characters.
- {Date(int)}: Date, e.g.: {Date(-3)}, generates a date 3 days ago, 2022-01-02.
- {Month(int)}: Date, e.g.: {Month(-3)}, generates the month 3 months ago, 2022-01.
- {Year(int)}: Date, e.g.: {Year(-3)}, generates the year 3 years ago, 2022.
- {IntStr(12)}: Generates a numeric string of length 12.
- {UpperStr(12)}: Generates an all-uppercase letter string of length 12.
- {LowerStr(12)}: Generates an all-lowercase letter string of length 12.
- {UpperLowerStr(12)}: Generates a string containing both uppercase and lowercase letters of length 12.
- {Str(12)}: Generates a string containing both uppercase and lowercase letters and numbers of length 12.
- {TimeGMT(-2)}: GMT date, e.g.: {TimeGMT(-3)}, generates a GMT date 3 days ago, Fri Aug 26 2022 23:59:59 GMT+0800.
- {DayBegin(1)}: Date, e.g.: {DayBegin(-3)}, generates the beginning date-time of 3 days ago, 2006-01-02 00:00:00.
- {DayEnd(-1)}: Date, e.g.: {DayEnd(3)}, generates the ending date-time of 3 days ahead, 2006-01-02 23:59:59.

#### Feature Data Auto-generation:
- {Date}: Generates the current date in the format 2006-01-02.
- {Time}: Generates the current date with hours, minutes, and seconds in the format 2006-01-02 10:59:05.
- {DeviceType}: Generates a device type.
- {QQ}: Generates a QQ number.
- {Age}: Generates an age.
- {Diploma}: Generates an educational qualification.
- {IntStr10}: Generates a numeric string of length 10.
- {Int}: Generates an integer.
- {Timestamp}: Generates a timestamp.
- {TimeFormat(2006-01-02 03:04:05.000 PM Mon Jan)}: Generates a formatted time, e.g., 2021-06-25 10:59:05.410 AM Fri Jun.
- {TimeFormat(15:04 2006/01/02)}: Generates a formatted time, e.g., 10:59 2021/06/25.
- {TimeFormat(2006-1-2 15:04:05.000)}: Generates a formatted time, e.g., 2021-6-25 10:59:05.410.
- {TimeFormat(XXX)}: Generates a formatted time based on custom requirements.
- {TimeGMT}: Generates a GMT formatted time, e.g., Fri Aug 26 2022 23:59:59 GMT+0800.
- {DayBegin}: Generates the start date-time of the current day, e.g., 2006-01-02 00:00:00.
- {DayEnd}: Generates the end date-time of the current day, e.g., 2006-01-02 23:59:59.
- {MonthBegin}: Generates the start date-time of the current month, e.g., 2006-01-01 00:00:00.
- {MonthEnd}: Generates the end date-time of the current month, e.g., 2006-01-31 23:59:59.
- {YearBegin}: Generates the start date-time of the current year, e.g., 2006-01-01 00:00:00.
- {YearEnd}: Generates the end date-time of the current year, e.g., 2006-12-31 23:59:59.
- {DayStampBegin}: Generates a timestamp for the start of the current day.
- {DayStampEnd}: Generates a timestamp for the end of the current day.
- {MonthStampBegin}: Generates a timestamp for the start of the current month.
- {MonthStampEnd}: Generates a timestamp for the end of the current month.
- {YearStampBegin}: Generates a timestamp for the start of the current year.
- {YearStampEnd}: Generates a timestamp for the end of the current year.
- {Int(100000,999999)}{Year(-58)}{Int(10,13)}{Int(10,28)}{Int(1000,9999)}: Generates combined data, such as a 58-year-old ID number.

#### Feature Data Custom Definition:
- JSON format: e.g., ```{"default": "v1,v2,v3,v4,……", "ch": "v1,v2,v3,v4,……", "en": "v1,v2,v3,v4,……"}.```
- Regular format: e.g., ```v1,v2,v3,v4,…….```

#### Cascading Data Custom Definition:
Note: (Data definitions are first made in system parameters, then accessed using {TreeData_Name}, the format refers to the China variable, currently supports only 3 levels, contact the administrator for more support if needed.)
- {TreeData_China[1]}: Generates cascading data, e.g., country-China.
- {TreeData_China[2]}: Generates cascading data, e.g., province-Zhejiang Province.
- {TreeData_China[3]}: Generates cascading data, e.g., city-Hangzhou City.