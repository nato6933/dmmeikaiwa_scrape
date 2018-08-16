# dmmeikaiwa_scrape

## Features
This command(not daemon) get DMM Eikaiwa teachers lesson schedule and will notify them by LINE.

- To filter teachers lesson schedule with setting time that you want.
- To customize notification template

Note:
    If there are no difference between previous data and present one, this command will not notify.

## Install & Quick start

1. Access below URL to get line notify token

https://notify-bot.line.me/my/

2. Go get and build

```bash
go get github.com/nato6933/dmmeikaiwa_scrape
cd dmmeikaiwa_scrape/
go build
```

3. Edit conf/setting.yaml

```bash
cp conf/setting.yaml.sample conf/setting.yaml
vim conf/setting.yaml
```

Write your line access token, log directory path, start_time, end_time, and teachers' IDs you want to book.

Note:

    You will get ID from teachers personal page's URL.

    e.g. ID is 25711 in https://eikaiwa.dmm.com/teacher/index/25711/.

4. Try to use
```bash
./dmmeikaiwa_scrape
```

