# Messenger Xelvra: Digital Freedom Manifest

**Messenger Xelvra** is a peer-to-peer (P2P) communication platform designed to restore privacy, security, and user control over digital communication. The project aims to create a secure, efficient, and decentralized platform that pushes the boundaries of P2P communication capabilities.

## I. Vision

### A. The Problem

Current centralized communication platforms threaten user privacy by collecting, analyzing, and monetizing their data. User messages pass through servers beyond their control, making them vulnerable to surveillance and censorship.

### B. The Solution

Messenger Xelvra addresses this problem by providing a platform for direct, uncensored, and independent communication. The platform is designed to return control over data to users and restore trust in digital communication with emphasis on **extreme speed, minimal resource consumption, and top-tier security and robustness.**

## II. Architecture

### A. Principles

Xelvra Messenger architecture is built on principles of privacy protection, security, and decentralization. Key principles include:

1. **P2P Communication:** Direct, encrypted communication between users without intermediaries.
2. **Hybrid P2P Model:** Strategic use of direct P2P connections with priority and relay services as fallback to ensure functionality, sustainability, and trust building, without compromising user privacy. Implementation of **parallel transports (QUIC + pre-initialized TCP)** for latency minimization and resilience maximization, and **automated ICE framework with AI-driven prediction.**
3. **User Experience:** Multi-platform support, efficiency, and intuitive design with **aggressive optimization for low resource consumption, including Progressive Onboarding (visual P2P explanation and interactive demo with local network simulator) and full Accessibility (WCAG 2.1 AA compliance, screen reader support).**

### **B. Technická architektura**

Messenger Xelvra je modulární systém skládající se ze tří hlavních komponent:

1. **peerchat-cli (Go):** Nástroj příkazové řádky pro vývoj a testování P2P logiky a pro spouštění P2P uzlu na pozadí jako systémová služba.  
2. **peerchat-api (Go):** Lokální API služba (gRPC) pro komunikaci s frontendovými aplikacemi, poskytující robustní rozhraní k P2P jádru a využívající **event-driven architekturu.**  
3. **peerchat\_gui (Flutter):** Multiplatformní grafické uživatelské rozhraní, optimalizované pro mobilní zařízení s důrazem na energetickou efektivitu.

![][image1]

**Detailní popis Go modulů (peerchat/internal/):**

* **p2p/:** Správa P2P sítě (go-libp2p, **jeden Kademlia DHT s lokální in-memory LRU caching vrstvou a adaptivním pollingem**, NAT traversal s **agresivním hole-punchingem, embedded STUN/TURN, AI-driven predikcí a port-knockingtactic**, správa spojení s mechanismy zotavení a pre-warmed TCP spojeními, **QUIC transport s kernel-level/user-space UDP batchingem, hardwarovou akcelerací a dynamickým škálováním okna (BBR+Cubic), TCP fallback, Onion routing pro *všechna* metadata s více vrstvami šifrování**, **Bluetooth LE/Wi-Fi Direct jako fallback pro mesh sítě s chytrým řízením spotřeby**).  
* **crypto/:** Implementace šifrovacích protokolů (Signal Protocol, X3DH, Double Ratchet), bezpečná správa klíčů s **Memory Hardeningem (mlock(), canaries, memguard)**, **ochrana proti Replay/DoS útokům (pokročilé rate-limiting, token buckets)**, a **odolnost proti timing útokům.**  
* **user/:** Správa uživatelských identit (**DID formát did:xelvra:\<hash\>, ověřování podpisy Ed25519 (ZKP plánováno pro Epochu 4\)**), **implementace vyhledávání peerů podle DID, robustní blokování uživatelů (s šifrovanou černou listou v DHT), Sybil Resistance (dynamický Proof-of-Work, systém automatické důvěry, omezení "Ghost" kontaktů).**  
* **message/:** Správa zpráv a souborů (přenos, offline zprávy, pub/sub, **komplexní správa skupin s rolími a pozvánkami**). Optimalizovaný přenos velkých souborů pomocí chunkingu.  
* **api/:** Implementace gRPC serveru a API handlerů s **robustním error handlingem, validací vstupů a rate limitingem. Zahrnuje monitoring pro Prometheus/Grafana a distribuované trasování s OpenTelemetry.**  
* **db/:** Abstrakce databázových operací (**SQLite s WAL módem pro vysoký výkon, nízkou latenci a odolnost proti korupci, jeden šifrovaný soubor userdata.db na uživatele s automatickým checkpointingem WAL souboru**).  
* **util/:** Pomocné funkce (logování, metriky, validace).

### **B.1 Protokolové specifikace**

* **Message Framing:** Všechny zprávy a datové pakety budou strukturovány pomocí Google Protobuf pro efektivní serializaci a deserializaci, zajišťující kompaktnost a rychlost přenosu.  
* **Handshake sekvence:** Detailní diagram toku pro X3DH a Double Ratchet protokol bude k dispozici v samostatné dokumentaci, popisující přesnou sekvenci výměny klíčů a navázání šifrované relace.  
* **Onion routing pro metadata:** Implementace vrstveného šifrování pro metadata (např. IP adresy peerů v DHT dotazech, časová razítka) inspirovaná principy Onion/Garlic routingu, aby se ztížila analýza síťového grafu a určení reálného zdroje/cíle komunikace pro externí pozorovatele.  
* **Mesh Networking Protocol (Příklad)**:  
  // pkg/proto/mesh.proto  
  syntax \= "proto3";

  package xelvra.mesh;

  message MeshPacket {  
      bytes sender\_id \= 1;  
      bytes message\_id \= 2;  
      uint32 hop\_limit \= 3; // TTL pro flooding  
      oneof payload {  
          bytes raw\_payload \= 4; // Šifrovaný a Onion routovaný obsah  
          // Volitelně specifické typy pro debug  
      }  
  }

  * **Transportní vrstva pro BLE:** Využije GATT profil s MTU (Maximum Transmission Unit) 512B pro efektivní přenos dat.  
  * **Wi-Fi Direct:** Aktivace pouze při dostatečné úrovni baterie (\>50%) vzhledem k její vyšší spotřebě.


## **III. Zabezpečení**

### **A. Filozofie**

Zabezpečení je klíčovým principem Xelvra Messengeru. Platforma chrání uživatele před různými hrozbami, včetně pasivního odposlouchávání, aktivních útoků a cenzury s důrazem na **proaktivní obranu, minimalizaci odhalených informací a odolnost vůči pokročilým hrozbám.**

### **B. Kryptografické jádro**

* **E2EE:** End-to-end šifrování zpráv mezi odesílatelem a příjemcem pomocí Signal Protocolu.  
* **Kryptografické primitivy:** Standardizované a prověřené algoritmy (AES-256, Curve25519, SHA-256/SHA-3, HKDF), **optimalizované využitím hardwarové akcelerace (AES-NI).**  
* **Bezpečná správa klíčů:** Generování, ukládání a správa klíčů s důrazem na **ochranu v paměti (mlock(), canaries, memguard) a na disku (šifrované soubory SQLite).**  
* **Key Rotation (Zero-Touch):** Automatická rotace dlouhodobých klíčů každých **60 dní** s uživatelskou notifikací **48 hodin** před expirací. Během "grace period" 72 hodin se provádí paralelní šifrování starým i novým klíčem pro zajištění bezztrátového přechodu. Udržování historie klíčů pro dešifrování starších zpráv. Tento proces bude plně automatizovaný a transparentní pro uživatele.  
* **Integrita dat:** Digitální podpisy a hašování pro ověření původu a integrity zpráv.  
* **Zero-Knowledge Proof:** Implementace ZKP mechanismů pro **ověřování identity bez odhalení citlivých informací (plánováno pro Epochu 4, prozatím Ed25519 podpisy).**

### **C. Ochrana metadat**

Minimalizace metadat a decentralizovaná identifikace uživatelů. **Onion routing pro obfuscaci *všech* metadat s cílem ztížit analýzu síťového grafu.**

### **C.6 Forward Secrecy**

* **Rotace klíčů:** Automatická rotace klíčů pro Double Ratchet algoritmus každých 100 odeslaných zpráv nebo po 24 hodinách nečinnosti, čímž se minimalizuje množství dat šifrovaných jedním klíčem.  
* **Automatická invalidace:** Klíče budou automaticky invalidovány po 7 dnech neaktivity v konverzaci, což zajistí, že staré relace nebudou představovat dlouhodobé riziko.

### **D. Odolnost sítě**

Ochrana proti **Sybil útokům (s dynamickým Proof-of-Work pro nové DHT záznamy, systémem automatické důvěry a omezením nových kontaktů pro "Ghost" uživatele)**, **DoS útokům (s pokročilým rate-limiting a správou spojení na více úrovních)** a **Replay útokům (s timestampy a sekvenčními čísly).**

### **C.7 Protection against Advanced Threats**

* **Sybil Resistance:**  
  * **Požadavek na Proof-of-Work:** Pro přidání nových záznamů do DHT (např. nových uživatelských identit) bude vyžadován dynamický Proof-of-Work, jehož obtížnost se bude měnit na základě síťové zátěže, aby se ztížil DDoS útok (flooding PoW požadavky).  
  * **Omezení kontaktů pro "Ghost" uživatele:** Uživatelé ve statusu "Ghost" budou mít omezený počet nových kontaktů, které mohou iniciovat za 24 hodin (např. 3/den), aby se zabránilo spamování.  
  * **Automatická důvěra:** Noví uživatelé mohou komunikovat s 5 kontakty/den bez CAPTCHA. Po ověření (např. QR kód od existujícího a důvěryhodného kontaktu) limity zmizí.  
* **Quantum Resistance:**  
  * **Hybridní šifrování:** Pro dlouhodobou ochranu budou kombinovány současné (např. X25519) a post-kvantové (např. Kyber768) algoritmy pro handshake fázi a navázání sdílených tajemství.  
  * **Možnost migrace:** Architektura bude navržena tak, aby umožňovala budoucí migraci na čistě post-kvantová kryptografická schémata, jakmile budou standardizována a prověřena.  
  * 

### **E. Transparentnost**

Open-source kód a nezávislé bezpečnostní audity.

## **IV. Ekosystém**

### **A. Hash Tokeny (HT)**

Interní virtuální kredity pro odměňování přínosu a zajištění udržitelnosti. HT nemají finanční hodnotu mimo ekosystém Xelvra Messengeru.

### **B. Cesta důvěry**

Systém pro budování důvěry a reputace mezi uživateli. Statusy uživatelů: Duch, Uživatel, Architekt, Ambasador, Bůh.

### **C. Komunitní správa**

Dlouhodobý cíl: Decentralizovaná správa sítě komunitou (DAO).

## **V. Obchodní model**

Transparentní a komunitně orientované financování. Crowdfunding a prodej HT pro financování dalšího rozvoje.

## **VI. Kvantifikovatelné cíle a Energetická optimalizace**

### **A. Kvantifikovatelné cíle výkonu a zdrojů**

* **Latence P2P zprávy (jedna cesta):**  
  * \< 50 ms pro přímá spojení.  
  * \< 200 ms přes relé.  
  * **Maximální latence při zátěži:** \< 100ms při 100 zprávách/s.  
* **Spotřeba paměti (CLI/Backend idle):** \< 20 MB (Go runtime).  
* **Paměťový limit při aktivním použití:** \< 50MB (Go runtime).  
* **Spotřeba CPU (CLI/Backend idle):** \< 1%.  
* **Latence API volání (interní):** \< 10 ms.  
* **Propustnost API:** \> 1000 RPC/s pro základní operace.

### **B. Energetická optimalizace**

Optimalizace spotřeby energie pro Go backend a Flutter frontend s **explicitními cíli pro mobilní zařízení:**

* **Spotřeba energie (mobilní, idle, pozadí):** \< 15 mW.  
* **Spotřeba energie (mobilní, aktivní chat):** \< 100 mW.  
* **Energetická náročnost (mobil):** \< 5% baterie/hod při aktivním chatování.  
* Implementace inteligentních strategií spánku a probuzení s využitím platformně specifických mechanismů (WorkManager pro Android, Background Fetch/VOIP Push pro iOS).

### **VI.B Implementační strategie**

* **Batching operací:** Shlukování menších síťových požadavků nebo databázových zápisů do větších dávek pro snížení režie a optimalizaci spotřeby energie. **Včetně kernel-level QUIC batchingu pro Linux.**  
* **Adaptivní polling:** Dynamická úprava frekvence heartbeatů, DHT dotazů a dalších periodických síťových aktivit na základě stavu baterie (např. L3 dotazy 1x za 10 min při \<20% baterie) a aktivity uživatele, aby se minimalizovala zbytečná aktivita na pozadí.  
* **Event-Driven Architektura:** Nahrazení polling mechanismů pro komunikaci mezi GUI a API a v rámci P2P sítě push notifikacemi (gRPC streams, WebSockets) pro snížení zátěže na CPU a baterii.  
* **Hardwarová akcelerace:** Využití instrukcí AES-NI (pokud jsou k dispozici na hardware) a dalších specifických instrukčních sad pro kryptografické operace za účelem zvýšení výkonu a snížení spotřeby CPU.  
* **Battery-Aware GC:** Statické nastavení GOGC=30 \+ ballast alloc (např. 1GB dummy array) pro stabilitu paměti a snížení frekvence GC cyklů.  
* **Deep Sleep Mode:** Při nízké úrovni baterie (\<15%) deaktivujte DHT a přepněte na "mesh-only" režim (pouze mDNS/Bluetooth LE/Wi-Fi Direct) pro minimální spotřebu.  
  * **Řešení konfliktních scénářů:**  
    * **Příchozí hovor:** Využití "light push" notifikací (např. FCM s vysokou prioritou, ale minimálním payloadem) pro lokální wake-up P2P uzlu. Očekávaná spotřeba: \~0.2 mW.  
    * **Důležitá zpráva:** Zpráva uložená v lokální mesh síti (přes BLE/Wi-Fi Direct) nebo DHT bude notifikována až při probuzení uzlu z Deep Sleep módu (např. pravidelné synchronizační okno). Očekávaná spotřeba: \~0.1 mW.  
    * **Systémové aktualizace:** Synchronizace aktualizací databáze/aplikace v definovaných časových oknech (např. každých 6 hodin) během noci nebo při připojení k nabíječce. Očekávaná spotřeba: \~0.3 mW během synchronizace.  
    * **Periodický ping (BLE beaconing):** Pro udržení minimální konektivity a usnadnění probuzení uzlu, i při vypnutém WiFi/Bluetooth.

    

## **VII. Nasazení a operace**

Strategie nasazení, distribuce a údržby aplikace.

### **VII. Deployment Strategy**

* **Automatizované buildy:** Plně automatizovaná CI/CD pipeline, která zahrnuje kompilaci, testování a podepisování všech binárních souborů a balíčků digitálními podpisy (s využitím Sigstore/cosign).  
* **Reproducible builds:** Zajištění, že binární soubory jsou reprodukovatelné z daného zdrojového kódu, což umožňuje nezávislým stranám ověřit integritu a absenci neoprávněných změn.  
* **Delta Updates (pro GUI):** Integrace **bsdiff** pro efektivní doručování malých aktualizací (\<100KB) sítě, minimalizující velikost stahovaných dat (kritické pro mesh sítě).  
* **Podpora offline aktualizací:** Možnost doručování aktualizací i v lokální mesh síti bez přístupu k internetu, což zvyšuje robustnost a autonomii uživatelů.

## **VIII. Testování a Zajištění Kvality**

Důkladné testování pro zajištění spolehlivosti a bezpečnosti, včetně:

* Unit testy  
* Integrační testy  
* E2E testy  
* Výkonnostní a zátěžové testy proti definovaným metrikám.  
* **Chaos Engineering:** Simulace selhání sítě a uzlů (např. náhodné pády uzlů, ztráta paketů, zpoždění, simulace výpadků internetu v Docker testech) pro ověření odolnosti systému v nepředvídatelných podmínkách.  
  * **Integrace do CI pipeline/docker-compose.yml:**  
    \# Příklad pro docker-compose.yml pro síťový chaos  
    services:  
      network-chaos:  
        image: nicholasjackson/chaos-http  
        command: \-target p2p-network \-latency 100ms \-jitter 50ms \-loss 10%  
        \# Přidejte sítě a závislosti pro cílení na p2p-network

* **Fuzzing:** Testování robustnosti parsování protokolů a vstupů pomocí generování náhodných, potenciálně škodlivých dat (zejména **pro QUIC handshake protokoly a Protobuf zprávy**).  
  * **Nástroje:** go-fuzz pro Go moduly.  
    \# Příklad použití go-fuzz  
    go-fuzz \-bin=./message-fuzzer.zip \-workdir=/fuzz

* **Penetrační testování:** Využití externích nástrojů a technik pro penetrační testování (např. Nmap, Metasploit, OWASP ZAP – pro testování exposed API, pokud bude relevantní, jinak pro síťovou vrstvu) pro odhalení zranitelností.  
  * **Penetrační testování QUIC handshake:** Otestujte pomocí [QUIC-Intruder](https://github.com/vanhauser-thc/thc-quic-intruder).  
  * **Side-channel útoky:** Ověřte odolnost proti side-channel útokům pomocí [CacheScout](https://github.com/cachescout/cachescout) (pro ověření AES-NI implementace a dalších kryptografických operací).  
  * **Odolnost proti timing útokům:** Analyzujte a přidejte umělá, konstantní zpoždění v kryptografických operacích (např. porovnávání klíčů), aby se zabránilo timing útokům.  
* **Testování v reálných sítích:**  
  * **Veřejné WiFi s captive portály:** Testování připojitelnosti a funkčnosti P2P.  
  * **Restriktivní firewally:** Ověření schopnosti procházet restriktivní firewally (porty 80/443).  
  * **Mobilní sítě s častým handoverem:** Testování odolnosti spojení a P2P sítě.  
* **Energy Profiling:** Integrujte do CI pipeline:  
  \# Pro Linux backend  
  perf stat \-e power/energy-pkg/ ./peerchat-cli test \--duration 5m  
  \# Pro Android frontend  
  adb shell dumpsys batterystats \--enable full-wake-history  
  adb bugreport \> bugreport.txt \# Analyzujte v Battery Historian

## **IX. Jak přispívat**

Informace o tom, jak se zapojit do vývoje a komunity Xelvra Messengeru. Pro usnadnění prvního spuštění a vývoje bude k dispozici příkaz peerchat-cli setup a **Docker-based testovací prostředí.**

## **X. Kodex Chování (Code of Conduct): Vytváříme Respektující Komunitu**

Xelvra je komunita postavená na důvěře, otevřenosti a spolupráci. Pro zajištění bezpečného, přívětivého a inkluzivního prostředí pro všechny, zavedli jsme tento Kodex Chování. Platí pro všechny účastníky projektu.

### **A. Naše Hodnoty**

* Respekt, inkluzivita, otevřenost, spolupráce, bezpečí.

### **B. Očekávané Chování**

* Být vstřícný a trpělivý.  
* Používat vítací a inkluzivní jazyk.  
* Být ohleduplný, poskytovat konstruktivní kritiku.  
* Respektovat odlišné názory.  
* Respektovat soukromí a bezpečnost.  
* Přijímat zodpovědnost za své chyby.

### **C. Nepřípustné Chování**

* Obtěžování, diskriminace, osobní útoky, trolling, škodlivý kód, zveřejňování soukromých informací, nátlak/vyhrožování.

### **D. Vymáhání Kodexu**

Případy porušení Kodexu Chování se řeší spravedlivě a transparentně.

## **XI. Licencování**

Messenger Xelvra je licencován pod **GNU Affero General Public License v3.0 (AGPLv3)**.

## **XII. Glosář**

* **Kademlia DHT:** Distribuovaná hašovací tabulka používaná pro objevování peerů a ukládání dat. V Xelvra Messengeru bude využívána s lokální in-memory LRU caching vrstvou.  
* **HT (Hash Token):** Interní virtuální kredity v ekosystému Xelvra Messengeru, sloužící k odměňování uživatelů za přínos k síti (např. relayování zpráv, udržování DHT uzlu) a zajištění udržitelnosti. HT nemají finanční hodnotu mimo ekosystém.  
* **Progressive Onboarding:** Uživatelsky přívětivý proces prvního spuštění, který postupně vysvětluje koncepty P2P a provádí uživatele nastavením, včetně vizuálních simulací sítě.  
* **Zero-Touch Šifrování:** Automatická správa kryptografických klíčů bez nutnosti manuálního zásahu uživatele, včetně automatické rotace, "grace period" a notifikací.  
* **Kernel-Level QUIC Batching:** Optimalizace přenosu dat v QUIC protokolu, kde jsou malé pakety shlukovány a odesílány přímo z jádra operačního systému (např. pomocí SO\_ZEROCOPY a io\_uring na Linuxu) pro snížení režie a zlepšení propustnosti. Pro ostatní OS fallback na user-space batching.  
* **Deep Sleep Mode:** Energeticky úsporný režim pro mobilní aplikace, kdy jsou síťové aktivity minimalizovány a DHT je deaktivováno ve prospěch lokální mesh komunikace (mDNS, Bluetooth LE/Wi-Fi Direct) pro maximální úsporu baterie, s řešením konfliktů probuzení (light push, BLE beaconing).  
* **AI-Driven Prediction / AI-Based Routing:** Použití lehkých modelů strojového učení (např. ONNX Runtime) k dynamickému vyhodnocování síťových podmínek a výběru nejefektivnějšího transportního protokolu nebo cesty zprávy v reálném čase, s validací vstupů a sandboxingem modelu.  
* **SQLite s WAL (Write-Ahead Logging):** Databázový režim, který zlepšuje výkon a odolnost proti pádům databáze, minimalizuje fragmentaci a umožňuje efektivní checkpointing.  
* **Port-Knocking:** Technika pro otevírání portů na firewallu odesláním předdefinované sekvence paketů na uzavřené porty, což zvyšuje obtížnost pro útočníky.

## **XIII. Průvodce řešením běžných problémů (Troubleshooting Guide)**

| Problém | Diagnostika | Řešení |
| :---- | :---- | :---- |
| **NAT failure** | peerchat-cli status \--nat | Aplikace by se měla pokusit o automatické řešení pomocí hole-punchingu a STUN/TURN. Pokud selže, aplikace se přepne na relé. Ruční konfigurace portů (port forwarding) na routeru je poslední možností. Zobrazí se diagnostický overlay v GUI. |
| **Nízká rychlost zpráv** | peerchat-cli status \--latency | Aplikace by měla automaticky vybírat nejrychlejší transport (QUIC/TCP/relé) na základě AI predikce. Zkontrolujte kvalitu sítě, případně se zkuste připojit k jinému relé serveru. Zobrazí se diagnostický overlay v GUI. |
| **Offline fungování** | peerchat-cli status \--network | Ověřte, zda je aktivní mDNS a zda jsou v lokální síti další peery. Zkontrolujte, zda je povoleno Bluetooth/Wi-Fi Direct pro mesh komunikaci. |
| **Vysoká spotřeba CPU** | Monitorujte procesy OS (top, htop). | Zkontrolujte logy aplikace. Aplikace automaticky optimalizuje GC a polling. Může indikovat problém s konkrétním modulem, který je třeba profilovat a optimalizovat. |
| **Nejde se připojit** | peerchat-cli connect \<peer\_id\> \--verbose | Zkontrolujte dostupnost cílového peeru. Aplikace by měla automaticky procházet NAT (včetně port-knockingu). Ověřte, zda firewall neblokuje komunikaci. |
| **Zprávy nedochází** | Zkontrolujte historii klíčů. | Klíče se automaticky rotují a invalidují po 7 dnech neaktivity. Zkontrolujte logy na obou stranách pro detaily handshake a doručení. Využijte "grace period" pro synchronizaci klíčů. |
| **Problémy s ID** | peerchat-cli id | Zkuste zregenerovat ID (pozor na ztrátu historie). Ověřte, zda se ID správně publikuje do DHT. Aplikace má mechanismy Sybil resistance a automatické ověření důvěry. |
| **GUI se zpomaluje** | Sledujte FPS v nastavení GUI nebo pomocí Flutter DevTools. | Zkontrolujte RepaintBoundary a použití const widgetů. Optimalizujte ListView.builder s ItemExtent. Snižte frekvenci animací. |

[image1]: <data:image/png;base64,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>