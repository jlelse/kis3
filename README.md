# KISSS (Keep It Simple Stupid Stats) a.k.a kis3

> I started to develop KISSS because I need a simple and privacy respecting method to collect visitor statistics of my website. I don't need any fancy dashboard, but I want to get exactly the stats I need. I need something fast, that is able to run even on really low end hardware like a Raspberry Pi.

## How to install

KISSS is really easy to install via Docker.

    docker run -d --name kis3 -e ENVVAR=VARVALUE -e ... -p 8080:8080 -v kis3:/app/data kis3/kis3

Replace ENVVAR and VARVALUE with the environment variables from the configuration.

Depending on your setup, replace `-p 8080:8080` with your custom port configuration. KISSS listens to port 8080 by default, but you can also change this via the configuration.

To persist the data KISSS collects, you should mount a volume or a folder to `/app/data`. When mounting an folder, give writing permissions to UID 100, because it is a non-root image to make it more secure.

It's also possible to use KISSS without Docker, but for that you need to compile it yourself. In the future there will be executables without dependencies available.

## Configuration

You can configure some settings using environment variables:

`PORT` (`8080`): Set the port to which KISSS should listen

`DNT` (`true`): Set whether or not KISSS should respect Do-Not-Track headers some browsers send

`DB_PATH` (`data/kis3.db`): Set the path for the SQLite database (relative to the working directory - in the Docker container it's `/app`).

You can make the statistics private and only accessible with authentication by setting both `STATS_USERNAME` and `STATS_PASSWORD` to a username and password. If only one or none is set, the statistics are accessible without authorization and public to anyone.

## Add to website

You can add the KISSS tracker to any website by putting `<script src="https://yourkis3domain.tld/kis3.js"></script>` just before `</body>` in the HTML. Just replace `yourkis3domain.tld` with the correct address.

## Requesting stats

You can request statistics via the `/stats` endpoint and specifying filters via query parameters (`?view=hours&format=chart` etc.). By combining this filters, you can exactly request the stats you want to get.

The following filters are available:

`view`: specify what data (and it's view counts) gets displayed, you have the option between `pages` (tracked URLS), `referrers` (tracked refererrers - only hostnames e.g. google.com), `useragents` (tracked useragents with version - browsers or crawl bots with version), `useragentnames` (tracked useragents without version), `hours` / `days` / `weeks` / `months` (tracks grouped by hours / days / weeks / months), `allhours` / `alldays` (tracks grouped by hours / days including hours or days with zero visits, spanning from first to last track in selection)

`from`: start time of the selection in the format `YYYY-MM-DD HH:MM`, e.g. `2019-01` or `2019-01-01 01:00`

`to`: end time of the selection

`url`: filter URLs containing the string provided, so `word` filters out all URLs that don't contain `word`

`ref`: filter referrers containing the string provided, so `word` filters out all refferers that don't contain `word`

`ua`: filter user agents containing the string provided, so `Firefox` filters out all user agents that don't contain `Firefox`

`orderrow`: row to use for ordering, `first` for the data groups, `second` for the view counts

`order`: select whether to use ascending order (`ASC`) or descending order (`DESC`)

`format`: the format to represent the data, default is `plain` for a simple plain text list, `json` for a JSON response or `chart` for a chart generated with ChartJS in the browser
