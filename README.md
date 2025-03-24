# DoubleCannon

Projekt **DoubleCannon** to aplikacja napisana w języku Go z wykorzystaniem biblioteki raylib-go. Symuluje ona wystrzał pocisku o zmiennej trajektorii w interaktywnym środowisku 3D.  

Główne funkcje programu:  
- Możliwość regulowania siły wystrzału poprzez przytrzymanie lewego przycisku myszy.  
- Wizualizacja przewidywanej trajektorii pocisku przed jego wystrzałem.  
- Symulacja lotu dwóch pocisków: czerwonego (pierwszy wystrzelony) oraz zielonego (drugi, który podąża tą samą ścieżką).  

Drugi pocisk, wystrzelony po naciśnięciu spacji, dociera do celu jednocześnie z pierwszym, uderzając w ten sam punkt.

## Jak korzystać

1. **Kompilacja i uruchomienie:**
   - Upewnij się, że masz zainstalowanego Go oraz bibliotekę [raylib-go](https://github.com/gen2brain/raylib-go).
   - Skopiuj projekt do swojego środowiska.
   - Skorzystaj z polecenia:
     ```
     go run main.go
     ```
     
2. **Interakcja:**
   - Przytrzymaj lewy przycisk myszy, aby "nabijać" siłę wystrzału. Siła zwiększa się wraz z czasem przytrzymania przycisku.
   - Po zwolnieniu lewego przycisku, wystrzelony zostaje czerwony pocisk, którego trajektoria zostanie obliczona i narysowana.
   - Naciśnij przycisk spacji, aby uruchomić zielony pocisk, kontynuujący lot na podstawie pozostałego czasu lotu.
   - Sterowanie kamerą odbywa się klawiszami (m.in. lewy shift i lewy ctrl do zmiany wysokości kamery).
   - Naciśnij "P", aby zatrzymać lub wznowić symulację.

## Inne informacje

Projekt zawiera również losowo rozmieszczone przeszkody, które są wizualizowane jako sześciany w scenie 3D. Trajektoria pocisku jest wizualizowana zarówno za pomocą linii, jak i małych kul.