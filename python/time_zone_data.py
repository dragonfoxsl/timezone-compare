from datetime import datetime
import pytz

print('US TimeZones')
for timeZone in pytz.country_timezones['US']:
    print(timeZone)


print('All Timezones')
for timeZone in pytz.all_timezones:
    print(timeZone)
    

print('US Time')
aware_est = datetime.now(pytz.timezone('US/Eastern'))
print('US EST DateTime', aware_est)

aware_pst = datetime.now(pytz.timezone('US/Pacific'))
print('US PST DateTime', aware_pst)
