from pkg.infra.pkl.pkl_importer import PklImporter
from unittest.mock import patch, mock_open


class TestPklImporter:
    
    @patch('builtins.open', new_callable=mock_open, read_data='a\n b\n c')
    @patch('os.listdir')
    def test_import_pkl(self, mock_listdir, mock_open):

        expected_files = [
            "temp_20240101-aaa.pkl",
            "temp_20240102-bbb.pkl"
        ]
        mock_listdir.return_value = expected_files

        pkl_importer = PklImporter("path/to/pkl_directory")

        import_data = pkl_importer.import_data()
        result1 = next(import_data) 
        result2 = next(import_data)

        mock_listdir.assert_called_once_with("path/to/pkl_directory")
        assert result1 == 'a'
        assert result2 == 'b'


