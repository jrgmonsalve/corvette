# corvette
## load data to zincsearch

folder structure

>/enron_mail_20110402/maildir/{users}/{folders}/{subfolders}/{email_file}.

email file content
```
Message-ID: <10077309.1075842026239.JavaMail.evans@thyme>
Date: Wed, 2 Jan 2002 14:01:06 -0800 (PST)
From: dawn.doucet@enron.com
To: patricia.henry@enron.com
Subject: Tax Adjustment re Zufferli
Cc: john.zufferli@enron.com
Mime-Version: 1.0
Content-Type: text/plain; charset=us-ascii
Content-Transfer-Encoding: 7bit
Bcc: john.zufferli@enron.com
X-From: Doucet, Dawn </O=ENRON/OU=NA/CN=RECIPIENTS/CN=DDOUCET>
X-To: Henry, Patricia </O=ENRON/OU=NA/CN=RECIPIENTS/CN=Phenry2>
X-cc: Zufferli, John </O=ENRON/OU=NA/CN=RECIPIENTS/CN=Jzuffer>
X-bcc:
X-Folder: \ExMerge - Zufferli, John\Inbox
X-Origin: ZUFFERLI-J
X-FileName: john zufferli 6-26-02.PST

Tricia, John has a stock option amount reported on his final pay stub that needs to be adjusted. The amount represents taxes on stock options and restricted stock exercised in 2001.  The stock option numbers on the memos sent by PaineWebber are correct, however, the restricted stock exercise should not be reported on the T4 (as per verification by Dave Peters at PriceWaterhouseCoopers on Jan 2/02).  Please make the appropriate adjustments so that the T4 is correct and advise me when corrected.  Thanks!

```

## how to start zincsearch

https://zincsearch-docs.zinc.dev/quickstart/

aditional download the file https://github.com/zincsearch/zincsearch/blob/main/.env.sample and rename to .env


## how to run script
> go run .\cmd\cli\ C:\projects\truora\enron_mail_20110402\maildir

## how generate cpu and mem graphs

> go tool pprof -svg cpu_profile.prof > cpu_profile.svg

> go tool pprof -svg mem_profile.prof > mem_profile.svg


## mejoras
1. que si se detiene por alguna razon el proceso, pregunte si desea comenzar de nuevo o retomar donde se quedo