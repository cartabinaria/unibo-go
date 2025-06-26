// SPDX-FileCopyrightText: 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

package rubrica

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {
	// Mock server to simulate the Unibo website
	handler := http.NewServeMux()
	handler.HandleFunc("/rubrica", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("b_start:int") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" lang="it" xml:lang="it">
<head><meta name="viewport" content="initial-scale=1.0, user-scalable=yes, width=device-width, minimum-scale=1.0" /><meta name="format-detection" content="telephone=no" /><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" /><meta http-equiv="Content-Type" content="text/html;charset=utf-8" /><meta name="key" content="Rubrica, Strutture, Persone" /><meta name="BestBet" content="Rubrica, Strutture" /><meta name="Description" content="UniboRubrica, rubrica dell&#39;Università di Bologna. UniboRubrica permette di cercare i contatti delle persone e delle strutture di Ateneo." /><meta name="ROBOTS" content="NOINDEX,NOFOLLOW" /><link rel="stylesheet" type="text/css" href="/uniboweb/resources/UniboSearch/styles/us4.css?rev=1.2" /><link rel="stylesheet" type="text/css" href="/uniboweb/resources/UniboSearch/styles/responsive.css?rev=1.0" />
    <script type="text/javascript" src="/uniboweb/resources/commons/js/jquery.min.js"></script>
    <script type="text/javascript" src="/uniboweb/resources/commons/js/modernizr-custom.js"></script>
    <script type="text/javascript" src="/uniboweb/resources/commons/js/unibo-cookies.min.js?v=20241130063953"></script><title>
	UniboRubrica - Università di Bologna
</title></head>

<!--
Server: WCSP-TWEB-01
Cached At : 16:25:54
 -->


<body>
    <div id="page">

        <!--MAIN NAVIGATION-->




<div id="mainNavigation">


        <h1 class="hidden">Rubrica</h1>


    <div id="usBar">
        <ul>
            <li>

                <a id="Header1_WebLink" accesskey="r" href="/uniboweb/unibosearch/default.aspx">Web Unibo</a>
             </li>

             <li>

                    <strong>Rubrica</strong>


             </li>

             <li>

                <a id="Header1_MappeLink" accesskey="g" href="/uniboWeb/unibomappe/default.aspx">Mappe</a>
             </li>


        </ul>
    </div>


    <div id="usRelated">
        <ul><li>
            <a href="https://idp.unibo.it/adfs/ls/infoSSO.aspx" target="_blank" title="Informazioni sul Single Sign-On di Ateneo (apre in nuova finestra)"><img src="/UniboWeb/Resources/UniboSearch/images/ssologo18x18.png" alt="Logo Single Sign-On di Ateneo"></a>
            <strong></strong> <a id="Header1_LoginLink" accesskey="l" href="/uniboweb/SignIn.aspx?ReturnUrl=%2funiboweb%2funibosearch%2frubrica.aspx%3ftab%3dPersonePanel%26mode%3dpeople%26query%3d%252bnome%253aantonio%2b%252bcognome%253acorradi">Login</a></li>
            <li><a id="Header1_PortaleLink" href="https://www.unibo.it/it">Portale di Ateneo</a></li>
            <li><a id="Header1_GuidaLink" href="/uniboweb/unibosearch/Guida.aspx">?</a></li>



                <li><a id="Header1_LangEn" accesskey="e" href="/uniboweb/unibosearch/rubrica.aspx?tab=PersonePanel&amp;mode=people&amp;query=%2bnome%3aantonio+%2bcognome%3acorradi&amp;lang=en">EN</a></li>

        </ul>
    </div>


    <div class="clearBoth"></div>
</div>



        <!--HEADER-->
        <div id="header">

            <!--LOGO-->

            <div id="logo" style="margin:1em 0pt;" >



                <a id="HomeLink" href="/UniboWeb/UniboSearch/Default.aspx"><img src="../sites/unibosearch/images/side_logo.png" alt="UniboSearch" /></a>
            </div>


            <!--SEARCH FORM-->

            <div id="searchForm">

                <!--SWITCH-->
                <div id="switch">

                            <ul>
                                <li class="active"><span>Persone</span></li>
                                <li><a href="/uniboweb/unibosearch/rubrica.aspx?tab=StrutturePanel">Strutture</a></li>
                                <li><a href="/uniboweb/unibosearch/rubrica.aspx?tab=FullTextPanel">Ricerca libera</a></li>
                            </ul>

                </div>


                <div id="usrInput">
                    <form method="post" action="/uniboweb/unibosearch/rubrica.aspx?tab=PersonePanel&amp;mode=people&amp;query=%2bnome%3aantonio+%2bcognome%3acorradi" id="Form1" class="search">
<div class="aspNetHidden">
<input type="hidden" name="__VIEWSTATE" id="__VIEWSTATE" value="3A4rl0Ce4P1Xki1qjlE1NFhwk7DSn977ru0kE+ntPJVHVbby48X+ORPgzsMFecpC9Iy2SO4wNQk/myui/FkcZBKTyu3SVnP3sBv56gkGCC2r4XaiBX58MgzYgZlLv8hUA9QFwAWXacu4wu8ZyYlV2Ez4kJF1q+GVL6U7MXDSS9gUgc2imGm/QomSC6ZScfeKQwPCSx+vcfYIkVHYtW0KmSdH3elQUGWbdVMsPBFrWWRcObMdDvtp8UmJF8NXVYHMA7dyZfwa0WCGNWEaD83U94WYdEbi+N2uahN/jY6z12G6wVxonpg3Vc5jsUNnppRjKLtXKp5Ump3uXIecEirF/4jwjTOm4b5bpmg9AK9yO5V8j8pB85Jahpmain+gFCbs" />
</div>

<div class="aspNetHidden">

	<input type="hidden" name="__VIEWSTATEGENERATOR" id="__VIEWSTATEGENERATOR" value="82AC61B8" />
	<input type="hidden" name="__EVENTTARGET" id="__EVENTTARGET" value="" />
	<input type="hidden" name="__EVENTARGUMENT" id="__EVENTARGUMENT" value="" />
	<input type="hidden" name="__EVENTVALIDATION" id="__EVENTVALIDATION" value="srk3G9OYX90DodCSdrwgVCRc9a2FXJFxkOEhPl9YOcMkOoZ6Zap0hZ/E2k6W+U8HW/Mbx3JzkE3y+t0QbaRwX44j5YPZm2zaP7pO+xD4G39emWaV1bhlBlz+1y5fIx7ffEaLZUz9Bjeb7UMZGz8sTVUlDtD/ZEQGM4fJsJLQEmUYJMOdKQwCPca73Hg/zsLgIwDT0WOzGoyLfJb/ykqyphkNpuHJHWO5z0bcm1VBvn8=" />
</div>




<div class="filterbyrole">
    <div>
        <input id="PersoneForm1_AllRoles" type="radio" name="PersoneForm1$Who" value="AllRoles" checked="checked" />
        <label for="PersoneForm1_AllRoles">Tutto il personale</label>
    </div>
    <div>
        <input id="PersoneForm1_Teachers" type="radio" name="PersoneForm1$Who" value="Teachers" />
        <label for="PersoneForm1_Teachers">Docenti e ricercatori</label>
    </div>
    <div>
        <input id="PersoneForm1_Staff" type="radio" name="PersoneForm1$Who" value="Staff" />
        <label for="PersoneForm1_Staff">Tecnici-amministrativi</label>
    </div>
</div>

<table>
    <tbody>




        <tr>
            <td style="width:45%">
                <div onkeypress="javascript:return WebForm_FireDefaultButton(event, &#39;PersoneForm1_SearchButton&#39;)">

                    <label for="PersoneForm1_SurnameField" class="screenreaders-helpers">Cognome</label>
                    <div class="borderd">
                    <input name="PersoneForm1$SurnameField" type="text" value="corradi" id="PersoneForm1_SurnameField" oninput="javascript:$(&#39;#resetSurname&#39;).toggle(!($(this).val() == &#39;&#39;));" placeholder="Cognome" />
                    <span><a id="resetSurname" tabindex="-1" style='display:block' href="#" class="nobottom-border" title="Cancella" onclick="javascript:$('#PersoneForm1_SurnameField').val('');$(this).hide();"><div class="is-depletable"></div></a></span>
                    </div>

</div>
            </td>
            <td style="width:45%">
                <div onkeypress="javascript:return WebForm_FireDefaultButton(event, &#39;PersoneForm1_SearchButton&#39;)">

                    <label for="PersoneForm1_NameField" class="screenreaders-helpers">Nome</label>
                    <div class="borderd">
                    <input name="PersoneForm1$NameField" type="text" value="antonio" id="PersoneForm1_NameField" oninput="javascript:$(&#39;#resetName&#39;).toggle(!($(this).val() == &#39;&#39;));" placeholder="Nome" />
                    <span><a id="resetName" tabindex="-1" style='display:block' href="#" class="nobottom-border" title="Cancella" onclick="javascript:$('#PersoneForm1_NameField').val('');$(this).hide();"><div class="is-depletable"></div></a></span>
                    </div>

</div>
            </td>
            <td style="width:10%;vertical-align:bottom;"><input type="image" name="PersoneForm1$SearchButton" id="PersoneForm1_SearchButton" src="/UniboWeb/Resources/UniboSearch/images/searchButton.png" alt="Cerca" /></td>
        </tr>


    </tbody>
</table>

<div id="advancedSearchLink">
    <ul><li><a id="PersoneForm1_FormLink" href="/uniboweb/unibosearch/rubrica.aspx?tab=PersonePanel&amp;mode=advanced">Ricerca alfabetica</a></li></ul>
</div>


<script type="text/javascript">
//<![CDATA[
var theForm = document.forms['Form1'];
if (!theForm) {
    theForm = document.Form1;
}
function __doPostBack(eventTarget, eventArgument) {
    if (!theForm.onsubmit || (theForm.onsubmit() != false)) {
        theForm.__EVENTTARGET.value = eventTarget;
        theForm.__EVENTARGUMENT.value = eventArgument;
        theForm.submit();
    }
}
//]]>
</script>


<script src="/UniboWeb/WebResource.axd?d=pynGkmcFUV13He1Qd6_TZBclabRlTOgEgRuVrbXCtgvzr-3cX0rLF70bHHt-3cwTEIx7deLXA878GFIgRuwamg2&amp;t=638628279619783110" type="text/javascript"></script>
</form>
                </div>

            </div>



            <!-- SCHEDA -->



            <!--DIRECTORY RESULTS-->




    <!--INFO BAR-->
	<div id="dirInfoBar">
		<span id="PersoneResults_RisultatiRicerca"><strong>1</strong> risultati per <strong>antonio corradi</strong></span>
	</div>



    <div id="results">





<div class="pages" id="pages_before">
    <ul>
        <li class="currentPageNumber">Pagina: </li>

                <li><strong>1</strong></li>

        <li class="totalPageNumber"> &nbsp; di &nbsp; 1</li>
    </ul>
</div>



        <!-- Risultati da Persone -->

                <table class="contact vcard">
                    <tbody>
                        <tr>
                            <th class="uid">15896</th>
                            <td class="fn name">Corradi, Antonio</td>
                        </tr>
                        <tr class="role">
                            <th></th>
                            <td>
                                Professore Alma Mater<br />Professore a contratto a titolo gratuito
                            </td>
                        </tr>





                        <tr>
                            <th>e-mail</th>
                            <td><a class="email" title='Scrivi una mail' href="mailto:antonio.corradi@unibo.it">antonio.corradi@unibo.it</a></td>
                        </tr>


                                    <tr>
                                        <th>tel</th>
                                        <td class="tel">
                                            <span class="type" style="display: none;">work</span>
                                            <span class="type" style="display: none;">tel</span>
                                            <span class="value"><a class="phone" href="tel:+39 051 20 9 3083">+39 051 20 9 3083</a></span>
                                        </td>
                                    </tr>

                                <!-- Sito Web Docente -->

                                <tr>
	                                <th>web</th>
		                            <td>
		                                <a class="url" href='/sitoweb/antonio.corradi'>https://www.unibo.it/sitoweb/antonio.corradi</a>
		                            </td>
	                            </tr>


                                <tr>
                                    <th>vcard</th>
                                    <td>
                                        <a href='/UniboWeb/vcard.vcf?UPN=antonio.corradi@unibo.it' title='scarica la Vcard'>
                                            <img src="/UniboWeb/Resources/UniboSearch/images/icons/ico_vcf.gif" alt="Vcard"/>
                                        </a>
                                    </td>
                                </tr>

                    </tbody>
                </table>






<div class="pages" id="pages_after">
    <ul>
        <li class="currentPageNumber">Pagina: </li>

                <li><strong>1</strong></li>

        <li class="totalPageNumber"> &nbsp; di &nbsp; 1</li>
    </ul>
</div>


    </div>





        </div>


        <!--FOOTER-->
        <div id="footer">



<p>&copy; 2025 - Universit&agrave; di Bologna -

    <a target="_blank" style="color:#bbb" href="https://www.unibo.it/it/ateneo/privacy-e-note-legali/privacy/informative-sul-trattamento-dei-dati-personali">Privacy</a> |
    <a href="#" style="color:#bbb" data-cc-open="">Impostazioni Cookie</a>



        </div>
    </div>

    <!-- Piwik -->
    <script type="text/javascript">
        var _paq = _paq || [];
        _paq.push(['setCookiePath', '/UniboWeb/UniboSearch']);
        _paq.push(['trackPageView']);
        _paq.push(['enableLinkTracking']);
        (function () {
            var u = "https://analytics.unibo.it/";
            _paq.push(['setTrackerUrl', u + 'piwik.php']);
            _paq.push(['setSiteId', 10875]);
            var d = document, g = d.createElement('script'), s = d.getElementsByTagName('script')[0];
            g.type = 'text/javascript'; g.async = true; g.defer = true; g.src = u + 'piwik.js'; s.parentNode.insertBefore(g, s);
        })();
    </script>
    <noscript>
        <p><img src="//analytics.unibo.it/piwik.php?idsite=10875" style="border: 0;" alt="" /></p>
    </noscript>
    <!-- End Piwik Code -->

</body>
</html>

		`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	// Override the baseUrl to point to the mock server
	baseUrl = server.URL + "/rubrica?query="

	contatti, err := Search("antonio", "corradi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(contatti) != 1 {
		t.Fatalf("expected 1 contact, got %d", len(contatti))
	}

	t.Run("LastName", func(t *testing.T) {
		if contatti[0].LastName != "Corradi" {
			t.Fatalf("expected last name 'Corradi', got %s", contatti[0].LastName)
		}
	})

	t.Run("FirstName", func(t *testing.T) {
		if contatti[0].FirstName != "Antonio" {
			t.Fatalf("expected first name 'Antonio', got %s", contatti[0].FirstName)
		}
	})

	t.Run("Email", func(t *testing.T) {
		if contatti[0].Email != "antonio.corradi@unibo.it" {
			t.Fatalf("expected email 'antonio.corradi@unibo.it', got %s", contatti[0].Email)
		}
	})

}
