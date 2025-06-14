# **Projekt Xelvra: Kompletní Task List pro Messenger P2P \- TODO.md**

Tento dokument představuje podrobný a konsolidovaný seznam úkolů rozdělený do fází, které se budou striktně dodržovat pro implementaci projektu Xelvra. Každý úkol by měl být dokončen s důrazem na kvalitu, testování a sebekontrolu.  
Důležité pokyny:

* **Před každou fází a po každém významném milníku prověř, zda jsi dokončil všechny úkoly a zda kód odpovídá specifikaci v README.md.**  
* **Vždy prioritizuj dokončení aktuálního úkolu a jeho testování před přechodem na další.**  
* **Veškeré změny kódu musí být commitovány do Gitu.** Viz sekce "VII. Verzování a Git".  
* **Využij paměť serveru MCP** pro ukládání informací o dokončených úkolech, výsledcích testů a řešeních problémů.  
* **Pokud narazíš na nejasnost nebo problém, který nelze vyřešit dle plánu, upozorni mě.**

### STRICT EXECUTION PROTOCOL  
**PŘED PŘECHODEM NA JAKÝKOLIV ÚKOL VŽDY OVĚŘ:**  

1. **ADHERENCE CHECK**  
   - [ ] Ověřil jsi, že kód **přesně odpovídá specifikaci v README.md** v těchto oblastech:  
     - Architektura modulů (`internal/p2p`, `internal/crypto`)  
     - Kvantifikovatelné metriky (latence <50ms, paměť <20MB idle)  
     - Bezpečnostní standardy (memory hardening, key rotation)  
   - [ ] Prokázal jsi soulad spuštěním **validace architektury**:  
     ```bash 
     ./arch_validator.sh  # Vlastní skript kontrolující strukturu projektu
     ```

2. **TASK COMPLETION VERIFICATION**  
   - [ ] Dokončil jsi **všechny dílčí úkoly** v aktuální fázi (např. všechny checkboxy v Epoše 1.1)  
   - [ ] Provedl jsi **předepsané testy** s explicitním výstupem:  
     ```go
     go test -v -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep "total:" | awk '{print $3}'
     ```
     - **Minimální pokrytí:** 80% pro kritické moduly (DHT, crypto)  
     - **Metriky splněny:** Latence <200ms relay, CPU <1% idle  

3. **GIT SANITY CHECK**  
   - [ ] Commity obsahují **explicitní reference na TODO položky**:  
     ```bash
     git log -1 --pretty=%B | grep -E "\[ \] TODO-#[0-9]+"  # Příklad: "[ ] TODO-#42 Implement QUIC batching"
     ```
   - [ ] Žádné **uncommitted změny** po dokončení úkolu:  
     ```bash
     git status --porcelain | wc -l | grep 0
     ```

4. **MCP REPORTING**  
   - [ ] Uložil jsi do MCP tyto **verifikované artefakty**:  
     - Výstup testů (včetně metrik spotřeby CPU/paměti)  
     - Screenshot CLI validace (pokud relevantní)  
     - Hash commitu (`git rev-parse --short HEAD`)  
   - [ ] Aktualizoval jsi stav v MCP formátem:  
     ```json
     {"epoch":1, "task":"P2P-DISCOVERY", "status":"DONE", "metrics":{"latency_ms":45, "mem_mb":18.7}}
     ```

5. **PROBLEM ESCALATION PROTOCOL**  
   Pokud narazíš na **nesoulad mezi kódem a specifikací**:  
   - [ ] **Okamžitě zastav práci** na úkolu  
   - [ ] Spusť diagnostický skript:  
     ```bash
     ./triage.sh <problem-module>  # Generuje report s debug logy
     ```
   - [ ] Nahlas problém do MCP s **povinnými položkami**:  
     - ID TODO úkolu  
     - Konkrétní bod specifikace (s číslem řádku z README)  
     - Navrhované řešení (nebo 3 možnosti)  
     - Výstup `triage.sh`

**ZÁKAZY:**  
❌ Žádné "dočasné workaroundy" v rozporu se specifikací  
❌ Žádné přesuny mezi epochami bez 100% dokončení předchozí  
❌ Žádné commity bez reference na TODO ID v commit message

## **I. Globální nastavení a Verzování v Gitu**

**Cíl:** Připravit vývojové prostředí a nastavit základní strukturu projektu s verzováním v Gitu, včetně robustního testovacího prostředí.

* \[ \] **Initalizace Git repozitáře:**  
  * \[ \] V kořenovém adresáři projektu (xelvra/) inicializuj nový Git repozitář.  
  * \[ \] Vytvoř počáteční .gitignore soubor s běžnými položkami pro Go, Flutter a operační systém Linux (např. bin/, pkg/, .idea/, .vscode/, build/, .DS\_Store, \*.env, \*.local, \*-lock.json, atd.).  
  * \[ \] Vytvoř počáteční commit s těmito soubory a prázdnou adresářovou strukturou.  
* \[ \] **Nastavení adresářové struktury:**  
  * \[ \] Vytvoř hlavní adresáře: cmd/, internal/, pkg/, tests/, peerchat\_gui/.  
  * \[ \] Vytvoř podadresáře uvnitř cmd/: peerchat-cli/, peerchat-api/.  
  * \[ \] Vytvoř podadresáře uvnitř internal/: p2p/, crypto/, user/, message/, util/, db/.  
  * \[ \] Vytvoř podadresáře uvnitř pkg/: proto/.  
* \[ \] **Počáteční commit struktury:**  
  * \[ \] Commituj nově vytvořené prázdné adresáře (nebo s prázdnými .gitkeep soubory) a soubor .gitignore. Použij zprávu commitu "Initial project setup and directory structure".  
* \[ \] **Nastavení Docker-based testovacího prostředí:**  
  * \[ \] Vytvoř docker-compose.yml soubor a Dockerfile pro peerchat-cli, které umožní snadné spouštění více instancí pro testování P2P komunikace.  
  * \[ \] Zahrň základní bootstrap uzly do docker-compose.yml pro síťovou inicializaci (např. 3-5 bootstrap uzlů).  
  * \[ \] Nastav DNS resolver v Docker Compose pro interní komunikaci mezi kontejnery.  
  * \[ \] **Git commit:** "chore: Setup Docker-based test environment for multi-node simulation"

## **II. Epoch 1: CLI (peerchat-cli) – Kompletní Implementace**

**Cíl:** Implementovat plně funkční a robustní CLI messenger se všemi základními i pokročilými funkcemi pro komunikaci v peer-to-peer síti, včetně interaktivního chatování a možnosti běhu na pozadí. **Klíčový důraz na extrémní rychlost, minimální spotřebu zdrojů, robustní bezpečnost a odolnost v nestabilních síťových podmínkách.**

* \[ \] **Inicializace Go projektu:**  
  * \[ \] V kořenovém adresáři projektu inicializuj Go modul: go mod init github.com/Xelvra/peerchat.  
  * \[ \] Přidej Go moduly pro go-libp2p (včetně quic-go), cobra, viper (pro konfiguraci), logrus (pro logování), gocui (pro interaktivní CLI UI), **go-sqlite3 (pro databázi),** github.com/distatus/battery (pro stav baterie), **github.com/microsoft/onnxruntime-go (pro AI-driven predikci).**  
  * \[ \] **Zvaž a přidej moduly pro efektivní práci s pamětí a CPU, které mohou pomoci s nenáročností (např. runtime/pprof pro profilování, github.com/awnumar/memguard pro memory hardening).**  
  * \[ \] **Git commit:** "Go module initialization and core dependencies for CLI"  
* \[ \] **A. Základní P2P Komunikace a Core Funkcionalita**  
  * \[ \] **Implementace go-libp2p uzlu a Optimalizace:**  
    * \[ \] V internal/p2p/node.go vytvoř strukturu pro PeerChatNode a funkci pro její inicializaci (vytvoření go-libp2p uzlu).  
    * \[ \] Nakonfiguruj základní transporty (QUIC, TCP).  
    * \[ \] Implementuj graceful shutdown uzlu, který zajistí uzavření všech spojení a uvolnění zdrojů.  
    * \[ \] **Prioritizuj efektivní využití systémových zdrojů (CPU, paměť) od samého začátku. Implementuj monitorování a minimalizuj režii go-libp2p uzlu.**  
    * \[ \] **Cíl pro idle režim:** CPU \< 1%, Paměť \< 20MB (Go runtime).  
    * \[ \] **Git commit:** "feat: Basic go-libp2p node initialization and shutdown with resource optimization focus"  
  * \[ \] **QUIC Transport:**  
    * \[ \] V internal/p2p/node.go nakonfiguruj go-libp2p uzel tak, aby **primárně a efektivně** využíval QUIC transport (quic-go).  
    * \[ \] **Zajisti, že implementace QUIC maximalizuje rychlost a minimalizuje latenci (Cíl: latence přenosu zprávy \< 50ms pro přímá spojení).**  
    * \[ \] **Kernel-Level QUIC Batching (Linux only):** Pro Linux (kernel 5.4+) využijte SO\_ZEROCOPY a io\_uring pro batchování UDP paketů přímo v kernelu, což sníží režii aplikační vrstvy.  
    * \[ \] **User-Space Batching (non-Linux fallback):** Pro Windows/macOS a starší Linux jádra implementujte aplikační batching s SO\_REUSEPORT pro efektivní využití socketů.  
    * \[ \] **QUIC paměťová kontrola:** Nakonfigurujte quic-go s AllowConnectionWindowIncrease=false a explicitními limity pro striktní kontrolu využití paměti. Implementujte dynamické škálování okna (ConnectionWindow) na základě hybridního algoritmu **BBR \+ Cubic** (s fallbackem na fixní okno při \>5% ztrátě paketů) pro optimalizaci propustnosti při vysokém zatížení a minimalizaci drops.  
    * \[ \] **Git commit:** "feat: Integrate and optimize QUIC transport with kernel-level/user-space batching, dynamic window scaling (BBR+Cubic), and strict memory control"  
  * \[ \] **Správa připojení a Robustnost sítě (Hybridní P2P model):**  
    * \[ \] V internal/p2p/connection.go implementuj funkce pro navazování a udržování spojení s jinými peery.  
    * \[ \] Zahrň logiku pro automatické znovupřipojení a **inteligentní správu spojení, která minimalizuje režii a zajišťuje robustnost i v nestabilních síťových podmínkách.**  
    * \[ \] **Implementuj mechanismy pro detekci a zotavení se z přerušených spojení (např. periodické "keep-alive" zprávy, detekce neaktivních spojení a jejich uzavření).**  
    * \[ \] **Paralelní transporty:** Před-inicializujte TCP připojení při startu uzlu jako "hot backup" (neaktivní, ale okamžitě použitelné) a využijte libp2p.Multiplex pro simultánní QUIC/TCP připojení, aby se minimalizovalo zpoždění (\~200ms) při selhání primárního QUIC transportu. **Pro aktivní konverzace udržujte "teplá" TCP spojení v pozadí (neaktivní, ale okamžitě použitelná) pro snížení latence při přepínání.**  
    * \[ \] **Explicitně prioritizuj přímá P2P spojení. Použití relay serverů (p2p-circuit) jako poslední možnost při selhání přímého spojení.**  
    * \[ \] **Git commit:** "feat: Robust connection management with explicit P2P prioritization, parallel QUIC/TCP transports, pre-warmed connections, and resilience mechanisms"  
  * \[ \] **Implementace objevování peerů (Kademlia DHT, mDNS, Broadcast, BLE/Wi-Fi Direct) a Mesh sítě:**  
    * \[ \] V internal/p2p/discovery.go implementuj:  
      * \[ \] **Rychlý UDP broadcast:** Pro okamžité nalezení peerů v lokální LAN na lokální subnet (např. 242.0.0.0/8).  
      * \[ \] **Kademlia DHT s lokální caching vrstvou (in-memory LRU):** Implementuj jeden Kademlia DHT. Lokální cache výsledků DHT dotazů do in-memory LRU cache (např. BigCache) s TTL (např. 5 minut) pro rychlejší dotazy na často komunikované peery. Implementujte mechanismy pro invalidaci cache při změnách sítě. Omezte dotazy na 1x za minutu při \>50% baterie.  
      * \[ \] **Prioritizace známých peerů:** Ukládejte IP adresy a PeerID často komunikovaných peerů do lokální in-memory LRU cache s automatickým TTL obnovením.  
      * \[ \] **mDNS:** Pro rychlé a energeticky efektivní objevování v lokální síti. Optimalizuj pro minimální síťový broadcast a rychlé vyhledávání.  
      * \[ \] **Bluetooth LE (BLE) (pro mobilní zařízení):** Implementujte jako fallback transport pro textové zprávy v mesh sítích bez internetu. Pro BLE použijte GATT profil s MTU=512B.  
      * \[ \] **Wi-Fi Direct (pro mobilní zařízení):** Implementujte jako fallback transport pro souborové přenosy v mesh sítích. Aktivujte pouze při \>50% baterie (vysoká spotřeba).  
    * \[ \] **Optimalizuj dotazy DHT pro rychlou odezvu a efektivitu.** Implementuj strategie pro omezení počtu dotazů (např. 1x za minutu při \>50% baterie; méně při nízké baterii).  
    * \[ \] **Adaptivní polling režim:** Implementujte v internal/p2p/discovery.go dynamické úpravy frekvence DHT dotazů a mDNS broadcastů podle stavu baterie (např. při \<20% baterie: DHT dotazy 1x za 10 min, snížená frekvence mDNS) a aktivity uživatele.  
    * \[ \] **Integrace Battery-Aware API:** V internal/p2p/node.go (nebo jiném vhodném místě) integrujte knihovnu github.com/distatus/battery pro získávání informací o stavu baterie v Go.  
    * \[ \] **Zajisti, že mDNS je energeticky efektivní a spolehlivý pro lokální síťovou komunikaci (ad-hoc mesh síť v případě absence internetu).**  
    * \[ \] **Explicitně ověř funkčnost sítě v offline/lokálním režimu bez přístupu k internetu, simulujíc mesh síť s využitím Bluetooth LE nebo Wi-Fi Direct (pro mobilní zařízení) jako fallback transportů.**  
    * \[ \] **Git commit:** "feat: Optimized peer discovery with UDP broadcast, Kademlia DHT with in-memory LRU cache, mDNS, adaptive polling, battery awareness, and multi-transport mesh capabilities including BLE/Wi-Fi Direct"  
  * \[ \] **Implementace průchodu NAT (ICE, STUN, TURN, Hole Punching) pro spolehlivost:**  
    * \[ \] V internal/p2p/nat.go implementuj plně **automatizovaný ICE framework** s agresivním hole-punchingem a paralelními pokusy o UDP/TCP spojení.  
    * \[ \] **Embedded STUN/TURN:** Integrujte lightweight STUN server (a volitelně TURN) přímo do bootstrap uzlů, aby uživatelé za NATem mohli automaticky využívat tyto veřejné uzly pro zjištění své veřejné IP a typ NATu.  
    * \[ \] **Integrovat testy STUN serverů (např. podle tools.ietf.org/html/rfc5389) pro ověření funkčnosti a výběr nejlepšího serveru.**  
    * \[ \] **Port-Knocking:** Pro restriktivní firewally (např. blokující vše kromě HTTP/S) implementujte port-knocking na TCP/443 jako pre-step před QUIC handshake, aby se "otevřely" porty.  
    * \[ \] Při selhání přímého spojení **automaticky přepněte na relay (p2p-circuit) bez uživatelského zásahu.**  
    * \[ \] **AI-Driven Prediction:** Prozkoumejte a integrujte jednoduchý ML model (např. ONNX runtime pro Go) k předpovídání síťových podmínek a výběru optimálního transportu (QUIC/TCP/Relay) na základě historických dat o latenci a úspěšnosti spojení.  
      // internal/p2p/ai\_routing.go  
      package p2p

      import (  
      	"time"  
      	"github.com/microsoft/onnxruntime-go" // Použít pro inferenci ONNX modelu  
      	// Další potřebné importy  
      )

      // NetworkConditions definuje vstupní parametry pro AI model.  
      type NetworkConditions struct {  
      	Latency        time.Duration // Aktuální latence spojení  
      	PacketLoss     float64       // Procento ztrátovosti paketů  
      	ConnectionType string        // Typ připojení (WiFi/Cellular/Ethernet)  
      	BatteryLevel   int           // Úroveň baterie v procentech  
      	SignalStrength int           // Síla signálu (např. RSSI pro Wi-Fi/BLE)  
      }

      // PredictOptimalTransport provádí inferenci ONNX modelu pro výběr optimálního transportu.  
      // Model je trénován na datech o síťových podmínkách a úspěšnosti transportů.  
      // Vrátí doporučený transport ("QUIC", "TCP", "RELAY", "MESH\_BLE", "MESH\_WIFI\_DIRECT").  
      func PredictOptimalTransport(conditions NetworkConditions) (string, error) {  
      	// TODO: Načíst a inicializovat ONNX model.  
      	// TODO: Převeďte 'conditions' na vstupní tensor pro ONNX model.  
      	// TODO: Spusťte inferenci modelu.  
      	// Model Output (příklad): Pravděpodobnosti pro každý transport.  
      	// Příklad vstupních pravidel:  
      	// \- Při baterii \<20% preferovat UDP broadcast/mDNS/BLE (nízkoenergetické).  
      	// \- Při vysokém PacketLoss (\>5%) a vysoké latenci (\>200ms) preferovat TCP/Relay.  
      	// \- Při slabém signálu preferovat mesh (BLE/Wi-Fi Direct).

      	// Zde bude logika pro převod výstupu modelu na string transportu.  
      	// Zajištění sandboxingu modelu a validace vstupů (rozsah latence, packet loss).

      	// Placeholder pro demonstraci.  
      	if conditions.BatteryLevel \< 20 {  
      		return "MESH\_BLE", nil  
      	}  
      	if conditions.PacketLoss \> 0.05 || conditions.Latency \> 200\*time.Millisecond {  
      		return "RELAY", nil  
      	}  
      	return "QUIC", nil  
      }

    * \[ \] **Git commit:** "feat: Automated NAT traversal with aggressive hole-punching, embedded STUN/TURN, AI-driven transport prediction, port-knockingtactic, and seamless relay fallback"  
* \[ \] **B. Kryptografické Jádro a Identita (Důraz na Bezpečnost a Výkon)**  
  * \[ \] **Základní end-to-end šifrování (E2EE):**  
    * \[ \] V internal/crypto/signal.go implementuj jádro Signal Protocolu – konkrétně fáze X3DH handshake pro navázání sdíleného tajného klíče.  
    * \[ \] Zajisti, že dočasné soukromé klíče a mezilehlé kryptografické hodnoty jsou po použití okamžitě vynulovány z paměti.  
    * \[ \] Implementuj Double Ratchet algoritmus pro evoluci klíčů a forward secrecy.  
    * \[ \] Implementuj AES-256 GCM pro šifrování samotných zpráv a HMAC pro integritu.  
    * \[ \] **Pečlivě zvaž výběr a implementaci kryptografických primitiv s ohledem na jejich bezpečnostní prověřenost i výkonnostní charakteristiky.**  
    * \[ \] **Hardwarová akcelerace kryptografie:** Povolte AES-NI v Go (nastavením env GOAMD64=v3 při buildu). Offloadujte kryptografické operace na dedikované gorutiny s runtime.LockOSThread() pro minimalizaci blokování hlavního vlákna a maximalizaci výkonu.  
    * \[ \] **Odolnost proti timing útokům:** Implementujte umělá, konstantní zpoždění v kryptografických operacích (např. porovnávání klíčů, generování nonces), aby se zabránilo timing útokům, kde by útočník mohl odvodit informace z doby trvání operací.  
    * \[ \] **Git commit:** "feat: E2EE with Signal Protocol, focusing on security, performance, hardware acceleration, and timing attack resistance"  
  * \[ \] **Ochrana proti specifickým útokům:**  
    * \[ \] V internal/crypto/security.go implementuj ochranu proti **Replay útokům** (např. pomocí timestampů a sekvenčních čísel s omezeným časovým oknem a bloom filtry pro detekci duplikátů).  
    * \[ \] Rozšiř ochranu proti **DoS útokům** (např. pokročilé rate-limiting příchozích spojení a zpráv na aplikační vrstvě, použití connection manageru go-libp2p, token buckets nebo leaky buckets algoritmy).  
    * \[ \] **Git commit:** "feat: Implement replay and enhanced DoS attack protections"  
  * \[ \] **Onion routing pro metadata:**  
    * \[ \] V internal/p2p/onion\_routing.go prozkoumej a začni s implementací základních principů onion routingu pro obfuscaci *všech* metadat (včetně DHT dotazů, signalizace přítomnosti a dalších síťových operací).  
    * \[ \] **Cílem je ztížit analýzu síťového grafu a určení reálného zdroje/cíle komunikace pro externí pozorovatele. Implementujte minimálně 3 vrstvy šifrování pro metadata.**  
    * \[ \] **Git commit:** "feat: Initial implementation of multi-layered onion routing for all metadata obfuscation"  
  * \[ \] **Správa klíčů a Key Rotation (Zero-Touch):**  
    * \[ \] V internal/crypto/key\_manager.go implementuj bezpečné generování, ukládání (šifrované v SQLite) a načítání kryptografických klíčů.  
    * \[ \] **Memory Locking:** Použijte mlock() nebo podobné techniky pro uzamčení citlivých dat v paměti (klíče, nonces) a zabránění jejich swapování na disk.  
    * \[ \] **Integrace Memguard:** Integrujte knihovnu memguard ([github.com/awnumar/memguard](https://github.com/awnumar/memguard)) pro automatické mazání bufferů obsahujících citlivá data po jejich použití.  
    * \[ \] Implementace "canaries" (ochranných hodnot) pro detekci přetečení bufferů (buffer overflow) a narušení paměti.  
    * \[ \] **Automatická rotace dlouhodobých klíčů každých 60 dní.**  
    * \[ \] **"Grace period" pro Key Rotation:** Při rotaci klíče udržujte starý klíč aktivní po dobu 72 hodin s paralelním šifrováním zpráv novým i starým klíčem, aby se předešlo ztrátě zpráv při desynchronizaci klíčů.  
    * \[ \] **Uživatelská notifikace 48 hodin před expirací dlouhodobých klíčů** (pro minimalizaci rušení UX).  
    * \[ \] **Udržování historie klíčů pro dešifrování starších zpráv.**  
    * \[ \] **Optimalizuj operace s klíči pro minimální zátěž CPU.**  
    * \[ \] **Git commit:** "feat: Secure and optimized key management with key rotation, grace period, memory hardening, and memguard"  
  * \[ \] **Decentralizovaná identita (DID) a Sybil Resistance:**  
    * \[ \] V internal/user/identity.go definuj strukturu pro MessengerID odvozené z kryptografického páru klíčů ve formátu **did:xelvra:\<hash\>**.  
    * \[ \] Implementuj funkci pro generování nového MessengerID (pár klíčů).  
    * \[ \] **Pro ověření identity použijte jednoduché podpisy Ed25519 (ZKP implementace je odložena do Epochy 4).**  
    * \[ \] Integruj mechanismus pro ukládání a zveřejňování veřejných klíčů prostřednictvím Kademlia DHT.  
    * \[ \] **Implementace vyhledávání peerů podle MessengerID (DID) v DHT.**  
    * \[ \] **Automatická důvěra a Sybil Resistance:** Noví uživatelé mohou komunikovat s 5 kontakty/den bez CAPTCHA. Po ověření (např. QR kód od existujícího a důvěryhodného kontaktu) limity zmizí.  
      * \[ \] **Proof-of-Work pro DHT záznamy:** Implementujte dynamickou obtížnost PoW pro nové DHT záznamy na základě síťové zátěže, aby se předešlo DDoS útokům (flooding PoW požadavky).  
    * \[ \] **Git commit:** "feat: Decentralized identity (DID:xelvra:) with Ed25519 signatures, DID lookup, and advanced Sybil resistance with automatic trust and dynamic PoW"  
* \[ \] **C. Zprávové a Souborové Služby**  
  * \[ \] **Systém zpráv:**  
    * \[ \] V internal/message/manager.go implementuj zpracování příchozích a odchozích zpráv.  
    * \[ \] Zahrň podporu pro offline zprávy (dočasné ukládání na DHT nebo přes relay, s vymazáním po doručení).  
    * \[ \] Implementuj Pub/Sub model pro skupinové chaty (využij go-libp2p PubSub).  
    * \[ \] **Správa skupin:** Detaily pro správu skupin (vytváření, pozvánky, přijímání/odmítání, opuštění, změna názvu, správa rolí členů) v internal/message/group\_chat.go.  
    * \[ \] **Git commit:** "feat: Message processing, offline messages, Pub/Sub, and comprehensive group management"  
  * \[ \] **Přenos souborů:**  
    * \[ \] V internal/message/file\_transfer.go implementuj bezpečný P2P přenos souborů (end-to-end šifrovaný).  
    * \[ \] Zahrň progress bar pro CLI a mechanismus pro obnovení přerušeného přenosu.  
    * \[ \] **Optimalizujte přenos velkých souborů pomocí chunkingu a paralelního streamování dat.**  
    * \[ \] **Git commit:** "feat: Secure P2P file transfer with progress, resume, and large file optimization"  
  * \[ \] **Blokování uživatelů:**  
    * \[ \] V internal/user/blocking.go implementuj funkcionalitu pro blokování nežádoucích kontaktů, včetně perzistentního ukládání seznamu blokovaných v SQLite.  
    * \[ \] Zajistěte, že zprávy od blokovaných uživatelů nejsou zobrazovány a pokusy o navázání spojení jsou odmítnuty na úrovni P2P uzlu.  
    * \[ \] **Šifrovaná černá lista v DHT s podpisy:** Pro DHT dotazy implementujte mechanismus pro šifrování informací o blokovaných uživatelích a jejich publikaci do DHT (s kryptografickým podpisem), aby se peerové mohli vyhnout zbytečným dotazům na blokované uzly.  
    * \[ \] **Git commit:** "feat: Implement persistent user blocking functionality with encrypted blacklist in DHT"  
* \[ \] **D. CLI Aplikace (peerchat-cli)**  
  * \[ \] **Základní CLI struktura (Cobra):**  
    * \[ \] V cmd/peerchat-cli/main.go vytvoř hlavní Cobra příkazy.  
    * \[ \] Nastav zpracování konfiguračních souborů pomocí viper (např. \~/.config/xelvra/peerchat-cli.yaml pro ID, klíče, nastavení DHT bootstrapů atd.).  
    * \[ \] Implementuj základní logování do souboru a na konzoli pomocí logrus s konfigurovatelnými úrovněmi logování.  
    * \[ \] **Git commit:** "feat: Cobra CLI structure, Viper config, and Logrus integration"  
  * \[ \] **Hlavní příkazy CLI:**  
    * \[ \] peerchat-cli init: Vygenerování nového MessengerID a inicializace konfiguračního souboru.  
    * \[ \] **peerchat-cli setup:** Nový příkaz pro první spuštění, který provede inicializaci ID, základní konfiguraci a ověří připojení k síti. Poskytne jasné pokyny pro uživatele. (Základ Progressive Onboarding pro CLI).  
    * \[ \] **peerchat-cli doctor \--fix:** Nový příkaz pro automatickou diagnostiku a pokus o opravu problémů s NAT traversalem a připojením, včetně integrace testů STUN serverů a kontroly firewallu. **Při detekci problémů (vysoká latence, časté relay připojení) automaticky spustit diagnostiku na pozadí a navrhnout opravu, případně ji automaticky provést s uživatelským souhlasem.**  
    * \[ \] peerchat-cli start: Spuštění P2P uzlu na pozadí (démonizace, pokud OS podporuje, nebo prostě jako dlouho běžící proces). Implementujte mechanismus pro spuštění jako systémová služba (systemd unit/launchd plist pro Linux/macOS).  
    * \[ \] peerchat-cli stop: Zastavení běžícího P2P uzlu.  
    * \[ \] peerchat-cli connect \<peer\_id\>: Pokus o navázání přímého P2P spojení.  
    * \[ \] peerchat-cli send \<peer\_id\> \<message\>: Odeslání E2E šifrované zprávy vybranému peeru.  
    * \[ \] peerchat-cli send-file \<peer\_id\> \<file\_path\>: Odeslání souboru.  
    * \[ \] peerchat-cli listen: Spuštění naslouchání na příchozí zprávy a zobrazení v konzoli (pro jednoduché testování).  
    * \[ \] peerchat-cli status: Zobrazení stavu P2P připojení, známých peerů, přijatých a odeslaných zpráv, stavu NAT, a energetického stavu (pokud je battery-aware API integrováno).  
    * \[ \] peerchat-cli id: Zobrazení vlastního MessengerID.  
    * \[ \] peerchat-cli discover: Ruční spuštění procesu objevování peerů a zobrazení nalezených.  
    * \[ \] peerchat-cli profile \<peer\_id\>: Zobrazení základních informací o vzdáleném peeru.  
    * \[ \] peerchat-cli manual: Zobrazení nápovědy a popisu všech příkazů a jejich použití.  
    * \[ \] **Git commit:** "feat: Implement all core peerchat-cli commands, including setup, automated doctor, and system service integration"  
  * \[ \] **Interaktivní chatovací rozhraní (gocui):**  
    * \[ \] V cmd/peerchat-cli/chat\_ui.go implementuj interaktivní textové uživatelské rozhraní pro chatování pomocí gocui.  
    * \[ \] Zobrazení přijatých zpráv, okno pro psaní zpráv, seznam aktivních chatů/kontaktů (včetně skupinových chatů).  
    * \[ \] Podpora pro přepínání chatů (individuální/skupinové).  
    * \[ \] Zobrazení stavu spojení a síťové kvality (např. ikony podobné GUI) v UI.  
    * \[ \] **Implementujte základní uživatelské příkazy uvnitř chatu (např. /block, /unblock, /join \<group\_id\>, /create group, /status).**  
    * \[ \] **Git commit:** "feat: Interactive CLI chat UI with advanced features and in-chat commands"  
  * \[ \] **Persistentní úložiště pro CLI (SQLite s WAL mode):**  
    * \[ \] V internal/db/sqlite.go implementuj **SQLite s WAL (Write-Ahead Logging) módem** jako vysokovýkonné lokální úložiště, které minimalizuje fragmentaci a zajišťuje lepší stabilitu databáze.  
    * \[ \] Zajisti, že databáze je uložena v **jednom šifrovaném souboru na uživatele (userdata.db)**.  
    * \[ \] Specifikace výhod SQLite s WAL: **lepší konkurence, odolnost proti korupci při pádu, konzistentní výkon.**  
    * \[ \] Vytvoř Repository pattern pro abstrakci databázových operací (např. SaveMessage, LoadMessages, SaveContact, LoadContact, SaveGroup, LoadGroups).  
    * \[ \] Zajisti šifrování citlivých dat v databázi (např. pomocí klíče odvozeného z uživatelského hesla nebo klíče odvozeného od hlavního klíče uzamčeného v paměti).  
    * \[ \] **SQLite WAL Fragmentace:** Implementujte automatický checkpoint WAL souboru každých 1000 transakcí nebo 50MB velikosti, aby se minimalizoval růst a fragmentace \-wal souborů.  
    * \[ \] **Git commit:** "feat: SQLite with WAL mode for persistent storage with single encrypted user file, robust repository pattern, and WAL checkpointing"  
* \[ \] **E. Testování a Zajištění Kvality pro Epoch 1**  
  * \[ \] **Unit Testy:**  
    * \[ \] V tests/unit/p2p\_test.go napiš unit testy pro internal/p2p (inicializace uzlu, objevování (včetně UDP broadcast a caching), NAT, connection management, transporty).  
    * \[ \] V tests/unit/crypto\_test.go napiš unit testy pro internal/crypto (X3DH, Double Ratchet, klíčová správa, šifrování/dešifrování zpráv, Replay/DoS ochrana, timing attacks).  
    * \[ \] V tests/unit/user\_test.go napiš unit testy pro internal/user (generování identity, DHT ukládání, blokování, Sybil resistance).  
    * \[ \] V tests/unit/message\_test.go napiš unit testy pro internal/message (zpracování zpráv, přenos souborů, správa skupin).  
    * \[ \] V tests/unit/db\_test.go napiš unit testy pro internal/db (databázové operace SQLite).  
    * \[ \] V tests/unit/util\_test.go napiš unit testy pro internal/util (logging, pomocné funkce).  
    * \[ \] **Spusť všechny unit testy a ujisti se, že projdou.**  
    * \[ \] **Git commit:** "test: Implement comprehensive unit tests for core Go modules"  
  * \[ \] **Integrační Testy (zvýšená pozornost):**  
    * \[ \] V tests/integration/cli\_test.go napiš integrační testy simulující komunikaci mezi více peerchat-cli instancemi (spouštění procesů na pozadí, využití Docker testovacího prostředí).  
    * \[ \] Ověř navazování spojení, odesílání/příjem šifrovaných zpráv, skupinový chat, přenos souborů a ukládání historie.  
    * \[ \] **Rozšíření testů o scénáře s nestabilní sítí, přerušením internetu (simulace tc qdisc add dev eth0 root netem loss 100% v Dockeru) a následným fungováním v mesh režimu (lokální mDNS, Bluetooth LE/Wi-Fi Direct simulace).**  
    * \[ \] **Automatizované testy přechodu LAN mesh → globální síť a naopak.**  
    * \[ \] **Testuj chování aplikace při vysokém zatížení zprávami a soubory, s důrazem na rychlost a latenci.**  
    * \[ \] **Git commit:** "test: Robust integration tests including mesh network scenarios, internet outage simulation, and network transitions"  
  * \[ \] **Výkonnostní a Zátěžové Testy (hlubší analýza a kvantifikovatelné metriky):**  
    * \[ \] Vytvoř komplexní skripty pro zátěžové testy (simulace velkého počtu zpráv, mnoho peerů v síti, simulace dlouhotrvajících konverzací).  
    * \[ \] Monitoruj spotřebu CPU, paměti, síťového provozu a **spotřebu energie (simulovanou nebo měřenou na referenčních zařízeních, pokud je to možné již v CLI fázi).**  
    * \[ \] **Energetický benchmarking:** Přidejte perf stat \-e power/energy-pkg/ (pro Linux) do CI pipeline pro energetické benchmarkování Go backendu.  
    * \[ \] **Reportujte mW/zprávu v testech výkonnosti.**  
    * \[ \] **Stanov a ověř cíle:**  
      * **Latence P2P zprávy (jedna cesta):** \< 50 ms pro přímá spojení, \< 200 ms přes relé.  
      * **Maximální latence při zátěži:** \< 100ms při 100 zprávách/s.  
      * **Spotřeba paměti (idle):** \< 20 MB (Go runtime).  
      * **Paměťový limit při aktivním použití:** \< 50MB (Go runtime).  
      * **Spotřeba CPU (idle):** \< 1%.  
      * **Energetická stopa (mobilní, idle):** \< 20 mW (odhad, upřesní se v Epoch 4).  
    * \[ \] **Analyzuj výsledky a proveď rozsáhlé optimalizace kódu pro rychlost a nenáročnost.**  
    * \[ \] **Git commit:** "perf: In-depth performance, load, and initial energy consumption tests with optimizations against quantifiable metrics and CI integration"  
  * \[ \] **Bezpečnostní Testy (interní a automatizované):**  
    * \[ \] Provedení interní revize kódu se zaměřením na bezpečnostní slabiny (např. úniky paměti, chyby v kryptografii, XSS v CLI, pokud relevantní).  
    * \[ \] **Integruj fuzzing (např. go-fuzz) pro testování robustnosti parsování protokolů a vstupů, se zaměřením na QUIC handshake pakety a Protobuf zprávy.**  
    * \[ \] **Zvaž implementaci nástrojů pro chaos engineering na lokální Docker síti (např. náhodné shazování uzlů, simulace ztráty paketů, zpoždění) pro ověření odolnosti.**  
      \# Example to add to docker-compose.yml for network chaos  
      services:  
        network-chaos:  
          image: nicholasjackson/chaos-http  
          command: \-target p2p-network \-latency 100ms \-jitter 50ms \-loss 10%  
          \# (Další konfigurace pro cílení na konkrétní síť Dockeru)

    * \[ \] **Git commit:** "fix: Address internal security audit findings, integrate fuzzing (including QUIC/Protobuf) and basic chaos engineering"  
  * \[ \] **Testování v reálných sítích:**  
    * \[ \] **Veřejné WiFi s captive portály:** Testování připojitelnosti a funkčnosti P2P v sítích s captive portály.  
    * \[ \] **Restriktivní firewally:** Ověření schopnosti aplikace procházet restriktivní firewally (např. s využitím portů 80/443 a relé serverů).  
    * \[ \] **Mobilní sítě s častým handoverem:** Testování odolnosti spojení a P2P sítě při častém přechodu mezi BTS (handover).  
    * \[ \] **Git commit:** "test: Implement real-world network scenario testing"  
  * \[ \] **Manuál a Dokumentace:**  
    * \[ \] Doplnění podrobného manuálu (man stránky nebo samostatný soubor MANUAL.md) pro peerchat-cli, popisující všechny příkazy, konfiguraci a řešení problémů.  
    * \[ \] Doplnění Godoc komentářů ke všem veřejným funkcím a strukturám v Go kódu.  
    * \[ \] **Git commit:** "docs: Complete CLI manual and GoDoc documentation"  
* \[ \] **F. Review a Finální Dokončení Epoch 1**  
  * \[ \] **Kódová revize (Code Review):**  
    * \[ \] Provedení finální kódové revize pro celou Epochu 1, aby se zajistilo dodržování kódovacích standardů, čistota kódu a absence "debtu".  
    * \[ \] **Git commit:** "refactor: Final code review and cleanups for Epoch 1"  
  * \[ \] **Testovací pokrytí:**  
    * \[ \] Zkontroluj testovací pokrytí (např. pomocí go test \-coverprofile=coverage.out).  
    * \[ \] Snaž se dosáhnout vysokého pokrytí (ideálně \>80%) pro kritické moduly.  
    * \[ \] **Git commit:** "test: Ensure high test coverage for Epoch 1"  
  * \[ \] **Finální ověření:** Ověření všech testů (unit, integračních) a ověření jejich úspěšného průchodu.  
  * \[ \] **Git commit:** "Epoch 1 (CLI) complete: All features implemented and tests passed"  
  * \[ \] **Uložení stavu epochy:** Ulož do MCP, že "Epoch 1 (CLI) je dokončena a otestována."


## **III. Epoch 2: API (peerchat-api) – Lokální Služba**

**Cíl:** Vytvořit robustní lokální API službu pro peerchat\_gui (a potenciálně další klienty), která bude efektivně komunikovat s P2P jádrem.

* \[ \] **Inicializace API projektu:**  
  * \[ \] V cmd/peerchat-api/main.go vytvoř základní strukturu pro Go API server.  
  * \[ \] Použij gRPC pro komunikaci s frontendovými aplikacemi (definicí .proto souborů v pkg/proto/).  
  * \[ \] Přidejte Go modul pro go-sqlite3 (pro databázi).  
  * \[ \] Integruj viper pro konfiguraci a logrus pro logování.  
  * \[ \] **Git commit:** "feat: Initial peerchat-api setup with gRPC"  
* \[ \] **Definice gRPC služeb:**  
  * \[ \] V pkg/proto/chat.proto definuj gRPC služby pro správu uživatelů, zpráv, kontaktů, skupin a stavu sítě.  
  * \[ \] Zahrň RPC metody pro odesílání/příjem zpráv (včetně streamování pro dlouhotrvající chat), správu kontaktů, získávání historie chatu, správu profilu uživatele, správu skupin (vytváření, pozvánky, atd.).  
  * \[ \] Generuj Go kód z .proto souborů.  
    // pkg/proto/mesh.proto (Příklad pro mesh protokol)  
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

  * \[ \] **Git commit:** "feat: Define comprehensive gRPC services and generate Go code, including mesh protocol protobuf"  
* \[ \] **Implementace gRPC serveru a handlerů (Event-Driven):**  
  * \[ \] V internal/api/server.go implementuj gRPC server a jeho metody.  
  * \[ \] Integruj API handlery s P2P logikou z internal/p2p, internal/crypto, internal/message a internal/user.  
  * \[ \] Zajisti asynchronní zpracování zpráv a událostí z P2P sítě (např. pomocí Go kanálů a non-blocking operací).  
  * \[ \] **Nahraďte polling za Event-Driven Architekturu s gRPC streams (Server-Side Streaming RPC) pro API, které umožní push notifikace z backendu do GUI místo neustálého dotazování.**  
  * \[ \] **Implementujte robustní API error handling a validaci vstupů.**  
  * \[ \] **Implementujte základní rate limiting na API úrovni pro prevenci zneužití.**  
  * \[ \] **Git commit:** "feat: Implement gRPC server with robust error handling, input validation, rate limiting, and event-driven architecture via gRPC streams"  
* \[ \] **Správa dat (Persistentní úložiště pro API \- SQLite s WAL mode):**  
  * \[ \] V internal/db/sqlite.go rozšiř stávající SQLite implementaci pro ukládání dat relevantních pro API (např. rozšířená historie zpráv, uživatelské profily, nastavení API).  
  * \[ \] **Optimalizujte dotazy a operace s SQLite pro maximální výkon a minimalizaci latence API.**  
  * \[ \] **Git commit:** "feat: Extend SQLite for API data persistence with performance optimizations"  
* \[ \] **Monitoring a Telemetrie API:**  
  * \[ \] Integrujte do API metriky pro Prometheus a Grafana (např. počet volání RPC, latence volání, počet chyb).  
  * \[ \] Nastavte export metrik na /metrics endpoint.  
  * \[ \] Implementujte distribuované trasování (např. s OpenTelemetry) pro sledování toku požadavků napříč komponentami.  
  * \[ \] **Git commit:** "feat: Implement API monitoring with Prometheus/Grafana metrics and OpenTelemetry tracing"  
* \[ \] **Testování Epoch 2:**  
  * \[ \] Napiš unit testy pro internal/api handlery.  
  * \[ \] Napiš integrační testy pro gRPC server (simulace klientů, testování streamingu, testování edge-case chování).  
  * \[ \] **Kvantifikovatelné cíle pro API:**  
    * **Latence API volání:** \< 10 ms (interní, bez P2P sítě).  
    * **Propustnost:** \> 1000 RPC/s pro základní operace.  
  * \[ \] **Git commit:** "test: Add comprehensive unit and integration tests for peerchat-api with performance targets"  
* \[ \] **Finální ověření Epoch 2:** Ověření všech testů a potvrzení funkčnosti API.  
* \[ \] **Git commit:** "Epoch 2 (API) complete: gRPC service implemented and tested"  
* \[ \] **Uložení stavu epochy:** Ulož do MCP, že "Epoch 2 (API) je dokončena a otestována."

## **III. Epoch 3: GUI (peerchat\_gui) – Multiplatformní Klient**

**Cíl:** Vytvořit intuitivní a uživatelsky přívětivé multiplatformní grafické uživatelské rozhraní pro Messenger Xelvra. **Klíčový důraz na nenáročnost, rychlost a extrémní optimalizaci pro spotřebu energie na mobilních zařízeních, s vynikající uživatelskou zkušeností a přístupností.**

* \[ \] **Inicializace Flutter projektu:**  
  * \[ \] V adresáři peerchat\_gui/ inicializuj nový Flutter projekt.  
  * \[ \] Přidej závislosti pro gRPC klienta (grpc), stavovou správu (Riverpod), notifikace, a další klienti-side knihovny (např. pro UI komponenty, lokální storage, pokud je potřeba).  
  * \[ \] **Pečlivě zvaž výběr všech závislostí s ohledem na jejich dopad na velikost aplikace, paměť a výkon. Minimalizujte externí závislosti.**  
  * \[ \] **Git commit:** "feat: Initial peerchat\_gui Flutter project setup with dependency optimization focus"  
* \[ \] **Integrace s gRPC API (Event-Driven):**  
  * \[ \] Generuj Dart kód z .proto souborů (stejné jako pro Go).  
  * \[ \] Vytvoř gRPC klienta ve Flutter aplikaci pro komunikaci s peerchat-api (lokálně běžící službou).  
  * \[ \] **Implementujte příjem zpráv a událostí z API pomocí gRPC streams, což nahradí polling mechanismy a umožní real-time push notifikace do GUI.**  
  * \[ \] **Git commit:** "feat: gRPC client integration for peerchat\_gui with event-driven message streaming"  
* \[ \] **Návrh a implementace UI/UX (s důrazem na výkon a nenáročnost):**  
  * \[ \] Navrhni uživatelské rozhraní (obrazovky pro přihlášení/registrace, seznam kontaktů, chatovací okno, nastavení, atd.) podle Material Design (nebo Cupertino pro iOS) směrnic.  
  * \[ \] **Progressive Onboarding:** Implementujte průvodce prvním spuštěním s vizuálním vysvětlením P2P konceptů a interaktivním demem s lokálním simulátorem P2P sítě (spouští se při prvním spuštění) bez nutnosti vytváření účtu, aby se usnadnil onboarding nových uživatelů. Implementujte vizuální průvodce šifrováním (animace klíčů X3DH).  
  * \[ \] Implementuj klíčové UI komponenty (seznam chatů, okno chatu, vstupní pole, zobrazení souborů, správa skupin).  
  * \[ \] Zajisti citlivý design pro různé velikosti obrazovek a orientace (telefon, tablet, desktop).  
  * \[ \] **Flutter optimalizace:**  
    * \[ \] Omezení animací na max 30 fps v nastavení, s možností uživatelského vypnutí.  
    * \[ \] Důsledné použití const widgetů a RepaintBoundary pro statické prvky pro minimalizaci zbytečných rebuildů a redraws.  
    * \[ \] **Optimalizujte použití ListView.builder a Sliver\* s ItemExtent pro efektivní rendering dlouhých seznamů (chat history) a předvídatelné scrollování na starších zařízeních.**  
    * \[ \] **Vyhni se zbytečným animacím a komplexním přechodům, které by mohly zvyšovat spotřebu energie.**  
  * \[ \] **Accessibility:** Zajištění souladu s WCAG 2.1 AA standardy pro přístupnost.  
  * \[ \] **Podpora pro screen readery, režimy vysokého kontrastu a navigaci pomocí klávesnice.**  
  * \[ \] **Síťový status v UI:** Přidejte ikony kvality spojení (🌐/🟢/🔴) s měřením latence v reálném čase. Vytvořte tooltipy/popisky s vysvětlením technických problémů (např. "Vysoká latence: Zkuste Wi-Fi" nebo "Jste v lokální mesh síti"). Zahrňte **Diagnostický overlay s detaily (NAT typ, použitý transport, packet loss)**, přístupný z UI.  
  * \[ \] **AI-Based Routing (vizuální reprezentace):** V GUI zobrazujte doporučenou cestu zprávy nebo stav sítě na základě predikcí AI modelu pro demonstraci optimalizace.  
  * \[ \] **Git commit:** "feat: Optimized UI/UX implementation for performance, low resource usage, accessibility, progressive onboarding, network status display with diagnostic overlay, and AI routing visualization with Flutter best practices"  
* \[ \] **Správa stavu aplikace (energeticky efektivní):**  
  * \[ \] Implementuj efektivní správu stavu pro celou aplikaci pomocí Riverpod.  
  * \[ \] **Nahrazení setState() správou stavu přes Riverpod s selektivními rebuildy pro minimalizaci zbytečných operací.**  
  * \[ \] Zajištění reaktivního zobrazení dat z API.  
  * \[ \] **Optimalizujte, aby se data načítala a aktualizovala pouze tehdy, když je to nezbytně nutné (tzv. "lazy loading" a "on-demand updates"), aby se šetřila baterie.**  
  * \[ \] **Git commit:** "feat: Energy-efficient state management for peerchat\_gui using Riverpod for selective rebuilds and lazy loading"  
* \[ \] **Uživatelské funkce:**  
  * \[ \] **Registrace/Přihlášení:** Možnost vytvořit novou identitu nebo importovat existující (s podporou importu klíčů/seed frází).  
  * \[ \] **Správa kontaktů:** Přidávání (pomocí DID), mazání, blokování kontaktů. Zobrazení statusu důvěry (Ghost, User, atd.).  
  * \[ \] **Individuální chat:** Odesílání a příjem textových zpráv s E2E šifrováním, zobrazení stavu doručení/přečtení.  
  * \[ \] **Skupinový chat:** Vytváření a správa skupin, odesílání zpráv, správa členství (přidávání/odebírání).  
  * \[ \] **Přenos souborů:** Odesílání a příjem souborů s vizuální indikací průběhu a možností pozastavení/obnovení.  
  * \[ \] **Historie chatu:** Načítání a zobrazení historie zpráv (optimalizované pro dlouhé historie).  
  * \[ \] **Uživatelský profil:** Zobrazení a úprava vlastního profilu, správa klíčů (např. export veřejného klíče).  
  * \[ \] **Nastavení:** Správa nastavení aplikace (např. notifikace, témata, jazyk, nastavení soukromí, konfigurace síťových priorit).  
  * \[ \] **Git commit:** "feat: Implement all core user functionalities in peerchat\_gui with advanced features"  
* \[ \] **Integrace notifikací:**  
  * \[ \] Implementuj push notifikace pro mobilní platformy (Firebase Cloud Messaging pro Android, Apple Push Notification service pro iOS) a desktop (systémové notifikace).  
  * \[ \] Umožněte uživatelům konfigurovat preference notifikací (zvuk, vibrace, obsah).  
  * \[ \] **Zajistěte, že notifikace na pozadí neprobouzejí aplikaci zbytečně často a jsou energeticky efektivní. Využijte platformně specifické mechanismy pro běh na pozadí (WorkManager pro Android, Background Fetch/VOIP Push pro iOS).**  
  * \[ \] **Git commit:** "feat: Integrate platform-specific push notifications with energy efficiency and background execution"  
* \[ \] **Testování Epoch 3:**  
  * \[ \] Napiš widget testy pro UI komponenty (pokrytí \> 80%).  
  * \[ \] Napiš integrační testy pro interakci s gRPC API a databázi.  
  * \[ \] **End-to-End testy pro GUI:** Automatizované testy simulující kompletní uživatelské scénáře (např. registrace, chatování, přenos souborů).  
  * \[ \] **Testování energetické náročnosti:** Integrace nástrojů jako Android Battery Historian nebo Xcode Instruments (Energy Log) do testovacího procesu pro měření reálné spotřeby energie a ověřování cílů.  
  * \[ \] **Git commit:** "test: Add comprehensive widget, integration, E2E, and energy consumption tests for peerchat\_gui"  
* \[ \] **Finální ověření Epoch 3:** Ověření všech testů a potvrzení funkčnosti GUI.  
* \[ \] **Git commit:** "Epoch 3 (GUI) complete: Multiplatform client implemented and tested"  
* \[ \] **Uložení stavu epochy:** Ulož do MCP, že "Epoch 3 (GUI) je dokončena a otestována."

## **IV. Epoch 4: Energetická Optimalizace a Zajištění Důvěry**

**Cíl:** Dále optimalizovat energetickou efektivitu celého systému a plně rozvinout systém důvěry, včetně pokročilé kryptografie.

* \[ \] **Energetická optimalizace (komplexní a finální):**  
  * \[ \] **Go Backend:** Detailní profilování výkonu a spotřeby energie Go backendu (peerchat-cli, peerchat-api). Finální optimalizace síťového provozu (např. dávkování zpráv, komprese) a cyklů CPU (např. efektivnější algoritmy, caching, optimalizace garbage collection).  
  * \[ \] **Flutter Frontend:** Finální optimalizace Flutter UI renderingu, minimalizace redraws, správa zdrojů pro **maximální snížení spotřeby baterie na mobilních zařízeních.** Využití nástrojů jako Flutter DevTools pro analýzu výkonu a spotřeby.  
  * \[ \] **Spánek/Probuzení (Deep Sleep Mode):** Implementace inteligentních strategií spánku a probuzení pro P2P uzel a GUI (včetně využití platformních mechanismů jako WorkManager pro Android a Background Fetch/VOIP Push pro iOS), aby se minimalizovala aktivní spotřeba energie, když aplikace není v popředí.  
    * **Konfliktní scénáře Deep Sleep Mode a jejich řešení:**  
      * **Příchozí hovor:** Využití "light push" notifikací (např. FCM s vysokou prioritou, ale minimálním payloadem) pro lokální wake-up P2P uzlu. Očekávaná spotřeba: \~0.2 mW.  
      * **Důležitá zpráva:** Zpráva uložená v lokální mesh síti (přes BLE/Wi-Fi Direct) nebo DHT bude notifikována až při probuzení uzlu z Deep Sleep módu (např. pravidelné synchronizační okno). Očekávaná spotřeba: \~0.1 mW.  
      * **Systémové aktualizace:** Synchronizace aktualizací databáze/aplikace v definovaných časových oknech (např. každých 6 hodin) během noci nebo při připojení k nabíječce. Očekávaná spotřeba: \~0.3 mW během synchronizace.  
      * **Periodický ping (BLE beaconing):** Implementujte periodický ping přes BLE beaconing (1x/15min) pro udržení minimální konektivity a usnadnění probuzení uzlu, i při vypnutém WiFi/Bluetooth.  
    * **Při \<15% baterie deaktivujte DHT a přepněte na "mesh-only" režim (pouze mDNS/Bluetooth LE/Wi-Fi Direct) pro minimální spotřebu.**  
  * \[ \] **Kvantifikovatelné cíle pro energetickou optimalizaci (mobilní):**  
    * **Spotřeba energie (mobilní, idle, pozadí):** \< 15 mW.  
    * **Spotřeba energie (mobilní, aktivní chat):** \< 100 mW.  
    * **Energetická náročnost (mobil):** \< 5% baterie/hod při aktivním chatování.  
  * \[ \] **Go GC Tuning (Battery-Aware):**  
    * \[ \] Dynamicky upravujte GOGC (např. GOGC=20 při \<30% baterie a GOGC=50 při dostatečné baterii) pro snížení latence GC pauz a zlepšení plynulosti.  
    * \[ \] **Statický GOGC \+ Ballast Alloc:** Pro větší stabilitu zvažte statické nastavení GOGC=30 (nebo jiné optimální konstantní hodnoty) a použití ballast allocation (např. 1GB dummy array) pro stabilizaci paměti a snížení frekvence GC cyklů. Dynamické změny GOGC mohou být nestabilní při dlouhém běhu.  
  * \[ \] **Git commit:** "perf: Comprehensive energy optimization for Go backend and Flutter frontend with deep sleep mode, battery-aware GC (static+ballast), and platform-specific background execution"  
* \[ \] **Rozšíření Cesty důvěry:**  
  * \[ \] Implementace uživatelských statusů: Ghost, User, Architect, Ambassador, God.  
  * \[ \] Definice kritérií pro povýšení mezi statusy (např. doba online, počet ověřených spojení, přínos pro síť, účast v komunitě, ověřené příspěvky).  
  * \[ \] Vizuální indikace statusů v GUI.  
  * \[ \] **Zero-Knowledge Proof pro identitu (ZKP Light):** Implementujte Schnorr podpisy pro ověřování identity, které nabízejí silné soukromí s nižší režií než plné zk-SNARKs.  
  * \[ \] **Git commit:** "feat: Implement Trust Path levels and criteria, including ZKP Light (Schnorr signatures) for identity verification"  
* \[ \] **Hash Tokeny (HT) – Rozšířená implementace:**  
  * \[ \] V internal/user/hashtoken.go vytvoř logiku pro generování a správu Hash Tokenů.  
  * \[ \] Implementuj mechanismus pro odměňování uživatelů HT za aktivní účast v síti (např. relayování zpráv, udržování DHT uzlu, poskytování relé služeb).  
  * \[ \] **Implementujte mechanismy pro férovou distribuci a prevenci zneužití systému odměn.**  
  * \[ \] **Git commit:** "feat: Extended Hash Token generation and robust reward mechanism"  
* \[ \] **Finální ověření Epoch 4:** Ověření všech optimalizací a funkčnosti systému důvěry.  
* \[ \] **Git commit:** "Epoch 4 (Optimization & Trust) complete"  
* \[ \] **Uložení stavu epochy:** Ulož do MCP, že "Epoch 4 je dokončena a otestována."

## **V. Epoch 5: Decentralizovaná Správa a Udržitelnost**

**Cíl:** Položit základy pro decentralizovanou správu a dlouhodobou udržitelnost projektu, s důrazem na kvantovou odolnost a bezpečnost dodavatelského řetězce.

* \[ \] **DAO (Decentralizovaná Autonomní Organizace) \- První kroky:**  
  * \[ \] Návrh základní struktury DAO pro správu Xelvra Messengeru (např. na bázi smart kontraktů na lightweight blockchainu nebo distribuovaného konsensu).  
  * \[ \] Implementace jednoduchého hlasovacího mechanismu (např. off-chain s ověřováním pomocí HT nebo on-chain pro klíčová rozhodnutí).  
  * \[ \] **Definujte proces pro návrhy, diskuse a hlasování o klíčových změnách protokolu nebo distribuci prostředků.**  
  * \[ \] **Git commit:** "feat: Initial DAO structure, basic voting mechanism, and governance process definition"  
* \[ \] **Rozšířená distribuce HT:**  
  * \[ \] Implementace mechanismů pro distribuci HT za specifické přínosy (např. řešení bugů, vývoj nových funkcí, tvorba obsahu, provozování bootstrap uzlů, poskytování relé služeb).  
  * \[ \] **Založte transparentní peněženku/fond pro správu HT.**  
  * \[ \] **Git commit:** "feat: Extended HT distribution mechanisms and transparent fund management"  
* \[ \] **Kvantová odolnost (finální integrace):**  
  * \[ \] **Důkladný průzkum a výběr post-kvantových kryptografických algoritmů (PQC) pro dlouhodobou ochranu.**  
  * \[ \] Integrace vybraných PQC algoritmů (Kyber768, Dilithium) do Signal Protocolu a správy klíčů (hybridní schéma).  
  * \[ \] **Git commit:** "feat: Full integration of quantum-resistant cryptography (hybrid scheme)"  
* \[ \] **Supply Chain Security:**  
  * \[ \] **Sigstore/cosign:** Implementujte digitální podepisování všech vydaných binárek, image Dockeru a dalších artefaktů pomocí Sigstore/cosign.  
  * \[ \] **SBOM generování:** Automaticky generujte a zveřejňujte SBOM (Software Bill of Materials) pro všechny komponenty, aby byla zajištěna transparentnost a sledovatelnost dodavatelského řetězce.  
  * \[ \] **Git commit:** "security: Implement robust supply chain security with Sigstore/cosign and SBOM generation"  
* \[ \] **Finální dokumentace:**  
  * \[ \] Kompletní aktualizace veškeré dokumentace (README, manuály, architektonické diagramy, glosář, průvodce řešením problémů).  
  * \[ \] Vytvoření sekce "Jak přispívat" s jasnými pokyny a etickými zásadami.  
  * \[ \] **Git commit:** "docs: Final project documentation update and comprehensive contribution guide"  
* \[ \] **Závěrečné testování a kontrola:**  
  * \[ \] Závěrečné komplexní testování celé aplikace (E2E testy).  
  * \[ \] **Penetrační testování:** Využití externích nástrojů a firem pro profesionální penetrační testování a bezpečnostní audit celého systému.  
    * **Penetrační testování QUIC handshake:** Otestujte pomocí [QUIC-Intruder](https://github.com/vanhauser-thc/thc-quic-intruder).  
    * **Side-channel útoky:** Ověřte odolnost proti side-channel útokům pomocí [CacheScout](https://github.com/cachescout/cachescout) (pro ověření AES-NI implementace a dalších kryptografických operací).  
    * **Odolnost proti timing útokům:** Analyzujte a přidejte umělá, konstantní zpoždění v kryptografických operacích (např. porovnávání klíčů), aby se zabránilo timing útokům.  
  * \[ \] Kontrola licenčních hlaviček ve všech souborech a zajištění souladu s AGPLv3.  
  * \[ \] **Git commit:** "chore: Final E2E testing, professional penetration testing, and license verification"  
* \[ \] **Finální ověření projektu:** Ověření všech testů (unit, integračních, E2E) a ověření jejich úspěšného průchodu proti všem definovaným kvantifikovatelným metrikám.  
* \[ \] **Git commit:** "Project final review and all tests passed"  
* \[ \] **Uložení finálního stavu projektu:** Ulož do MCP, že "Projekt PeerChat je dokončen dle plánu."

## **VII. Verzování a Git**

**Cíl:** Zajistit konzistentní a efektivní verzování všech změn kódu v Gitu.

* \[ \] **Pravidlo pro commity:** Po dokončení každého smysluplného úkolu (i dílčího), který mění kód, vytvoř commit.  
* \[ \] **Zprávy commitů:** Používej jasné a popisné zprávy commitů, které stručně shrnují, co bylo změněno a proč. Formát: Typ: Popis změny (např. feat: Implement basic X3DH handshake, fix: Correct NAT traversal issue, refactor: Clean up P2P node initialization).  
* \[ \] **Větve (branches):** Pro komplexnější úkoly nebo experimentální funkce vždy vytvoř samostatnou větev.  
  * \[ \] Pro novou fázi nebo velký úkol: git checkout \-b feature/faze-\<cislo\>-\<popis\>  
  * \[ \] Po dokončení a otestování: git merge feature/faze-\<cislo\>-\<popis\> do hlavní větve (např. main nebo develop).  
* \[ \] **Před commitem:** Vždy zkontroluj stav (git status) a ujisti se, že commituješ jen relevantní změny. Použij git add \-p nebo git add \<file\> pro selektivní stageování.  
* \[ \] **Push na GitHub:** Pravidelně pushuj změny na https://github.com/Xelvra/peerchat.
